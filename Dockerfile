#
# Builds antas-natif / antas-natif-turbo against the native (CGO) PDFium
# build, selected via the `natif` build tag (vs. the webassembly variant).
# Uses the plain, non-V8 PDFium release from bblanchon/pdfium-binaries.
#
# Build:
#   docker build -t antas-natif .
#
# Override versions if needed:
#   docker build --build-arg PDFIUM_VERSION=chromium%2F7934 --build-arg GO_VERSION=1.23 -t antas-natif .

ARG GO_VERSION=1.25
ARG PDFIUM_VERSION=chromium%2F7934
ARG PDFIUM_ARCHIVE=pdfium-linux-x64.tgz

#############################################
# Stage 1: fetch native PDFium (plain, no V8)
#############################################
FROM debian:bookworm-slim AS pdfium

ARG PDFIUM_VERSION
ARG PDFIUM_ARCHIVE

RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates curl \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /pdfium
RUN curl -fsSL -o pdfium.tgz \
		"https://github.com/bblanchon/pdfium-binaries/releases/download/${PDFIUM_VERSION}/${PDFIUM_ARCHIVE}" \
	&& tar -xzf pdfium.tgz \
	&& rm pdfium.tgz

#############################################
# Stage 2: build the Go binaries with CGO
#############################################
FROM golang:${GO_VERSION}-bookworm AS builder

# gcc/build-essential for cgo, pkg-config so go-pdfium's cgo build finds
# PDFium, libturbojpeg-dev for the pdfium_use_turbojpeg build tag.
RUN apt-get update \
	&& apt-get install -y --no-install-recommends \
		pkg-config \
		libturbojpeg0-dev \
	&& rm -rf /var/lib/apt/lists/*

COPY --from=pdfium /pdfium/lib /opt/pdfium/lib
COPY --from=pdfium /pdfium/include /opt/pdfium/include

RUN mkdir -p /opt/pdfium/lib/pkgconfig \
	&& printf 'prefix=/opt/pdfium\nlibdir=${prefix}/lib\nincludedir=${prefix}/include\n\nName: PDFium\nDescription: PDFium\nVersion: 7934\nRequires:\n\nLibs: -L${libdir} -lpdfium\nCflags: -I${includedir}\n' \
		> /opt/pdfium/lib/pkgconfig/pdfium.pc

ENV CGO_ENABLED=1
ENV PKG_CONFIG_PATH=/opt/pdfium/lib/pkgconfig
ENV LD_LIBRARY_PATH=/opt/pdfium/lib

WORKDIR /src

# Cache module downloads separately from source changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -tags="natif pdfium_use_turbojpeg" -o /out/antas-natif .

#############################################
# Stage 3: minimal runtime image
#############################################
FROM debian:bookworm-slim

RUN apt-get update \
	&& apt-get install -y --no-install-recommends \
		ca-certificates \
		libturbojpeg0 \
	&& rm -rf /var/lib/apt/lists/*

# Native PDFium shared library, made discoverable to the dynamic linker.
COPY --from=pdfium /pdfium/lib/libpdfium.so /usr/local/lib/libpdfium.so
RUN ldconfig

COPY --from=builder /out/antas-natif /usr/local/bin/antas-natif

WORKDIR /work
ENTRYPOINT ["/usr/local/bin/antas-natif"]
