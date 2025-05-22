package engine

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Stockfish struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Scanner
}

func NewStockfish(path string) (*Stockfish, error) {
	cmd := exec.Command(path)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stdoutPipe)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Stockfish{
		cmd:    cmd,
		stdin:  stdin,
		stdout: scanner,
	}, nil
}

func (s *Stockfish) WriteLine(line string) error {
	_, err := s.stdin.Write([]byte(line + "\n"))
	return err
}

func (s *Stockfish) ReadLineUntil(expect string) string {
	for s.stdout.Scan() {
		text := s.stdout.Text()
		if strings.Contains(text, expect) {
			return text
		}
	}
	return ""
}

func (s *Stockfish) GetBestMove(fen string) (string, error) {
	s.WriteLine("uci")
	s.WriteLine("isready")
	s.ReadLineUntil("readyok")

	s.WriteLine("ucinewgame")
	s.WriteLine(fmt.Sprintf("position fen %s", fen))
	s.WriteLine("go depth 12") // bisa kamu sesuaikan (semakin tinggi, semakin pintar)
	line := s.ReadLineUntil("bestmove")

	if !strings.HasPrefix(line, "bestmove") {
		return "", fmt.Errorf("unexpected Stockfish output: %s", line)
	}

	parts := strings.Split(line, " ")
	if len(parts) >= 2 {
		return parts[1], nil
	}

	return "", fmt.Errorf("no move found")
}
