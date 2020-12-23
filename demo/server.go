package main

type server interface {
	ListenAndServe() error
}

func runServer() error {
	s := initServer("127.0.0.1:8999", nil)
	err := s.ListenAndServe()
	return err
}
