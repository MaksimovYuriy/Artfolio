package artwork

func extensionForFormat(format string) (string, bool) {
	switch format {
	case "jpeg":
		return ".jpg", true
	case "png":
		return ".png", true
	default:
		return "", false
	}
}
