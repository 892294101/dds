package dat

/*
func (w *ReadCache) Init(s *spfile.Spfile, log *logrus.Logger) error {
	trail := s.GetTrail()

	w.ProcName = *s.GetProcessName()

	home, err := utils.GetHomeDirectory()
	if err != nil {
		return err
	}

	dir := *trail.GetDir()
	ok := strings.HasPrefix(dir, "./")
	if ok {
		ind := strings.LastIndex(dir, "/")
		if ind == -1 {
			return errors.Errorf("Trail directory extraction error: %s", dir)
		}
		w.DatDir = path.Join(*home, dir[:ind])
		w.Prefix = dir[ind+1:]
	} else {
		ok := strings.HasPrefix(dir, "/")
		if ok {
			ind := strings.LastIndex(dir, "/")
			if ind == -1 {
				return errors.Errorf("Trail directory extraction error: %s", dir)
			}
			w.DatDir = dir[:ind]
			w.Prefix = dir[ind+1:]
		}
	}

	if len(w.DatDir) == 0 || len(w.Prefix) == 0 {
		return errors.Errorf("Failed to load trail directory when loading writer: %s", dir)
	}

	return nil
}
*/
