package render

type UploadPolicy struct {
	lastSheetID uint16
}

func (p *UploadPolicy) NoteSheet(sheetID uint16) bool {
	if p == nil {
		return true
	}
	if p.lastSheetID == sheetID {
		return false
	}
	p.lastSheetID = sheetID
	return true
}
