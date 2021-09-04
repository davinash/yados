package events

import (
	"fmt"
	"testing"
)

func func5() error {
	return fmt.Errorf("<- func5")
}
func func4() error {
	err := func5()
	if err != nil {
		return fmt.Errorf("<- func4 %w", err)
	}
	return nil
}

func func3() error {
	err := func4()
	if err != nil {
		return fmt.Errorf("<- func3 %w", err)
	}
	return nil
}

func func2() error {
	err := func3()
	if err != nil {
		return fmt.Errorf("<- func2 %w", err)
	}
	return nil
}

func func1() error {
	err := func2()
	if err != nil {
		return fmt.Errorf("func1 %w", err)
	}
	return nil
}

func TestGoErrors(t *testing.T) {
	err := func1()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

}
