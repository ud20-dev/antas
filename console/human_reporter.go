package console

func (HumanReporter) PageRendered(pageNum int, path string) {
	PrintWithStyle(InfoStyle, "Rendered page %d to %s", pageNum, path)
}

func (HumanReporter) Done(outDir string, pageCount int) {
	PrintWithStyle(SuccessStyle, "All pages rendered to %s (%d pages)", outDir, pageCount)
}

func (HumanReporter) Error(err error) {
	PrintWithStyle(ErrorStyle, "%s", err)
}