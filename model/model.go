package model

const (
    COMPILE = "COMPILE"
    TEST = "TEST"
)

type Work struct {
    Command string
    Args []string
}

type Node struct {
    Id int
    Busy bool
    Abilities map[string]struct{}
}

type WorkStatus struct {
    Done bool
    Error string
    Results []byte
}

func (node Node) CanCompile() bool {
    _, ok := node.Abilities[COMPILE]
    return ok
}

func (node Node) CanTest() bool {
    _, ok := node.Abilities[TEST]
    return ok
}

