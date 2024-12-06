package POOL

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type test struct {
	Name string
	Val  int
}

/*
	func TestInit(t *testing.T) {
		// Test Full Process of pool
		var p, err = NewPool(test{}, "TEST", resultOfTest, 20, false, 1*time.Second)
		if err != nil {
			t.Fatalf("Init failled  = %v", err)
		}
		p.Init()
		time.Sleep(1 * time.Second)
		err = p.Add("Var1", test{"Var1", 1}, 5)
		if err != nil {
			t.Fatalf("Init failled  = %v", err)
		}
		// Try to re add the same key
		time.Sleep(1 * time.Second)
		err = p.Add("Var1", test{"Var1", 1}, 5)
		if err != nil {
			t.Fatalf("Init failled  = %v", err)
		}
		// Try to  add new key
		time.Sleep(1 * time.Second)
		for j := 0; j < 3000; j++ {
			k := "VarJ" + fmt.Sprint(j)
			err = p.Add(k, test{k, 3 * j}, 1)
		}
		err = p.Add("Var2", test{"Var2", 2}, 1)
		if err != nil {
			t.Fatalf("Init failled  = %v", err)
		}
		time.Sleep(60 * time.Second)

}
*/
func resultOfTest(p *PoolEntry) error {
	t, ok := p.Content.(test)
	if ok {
		fmt.Println("JE demarre", p.Key, "avec une attente de ", t.Val)
	}

	time.Sleep(time.Duration(t.Val) * time.Millisecond)
	if t.Val > 170 {
		fmt.Println("Merde ca marcher pas")
		return errors.New("Ko")
	}
	fmt.Println("Ok ca marcher")
	return nil
}

func TestTriggerError(t *testing.T) {
	var p, err = NewPool(test{}, "TEST", resultOfTest, 20, false, 1*time.Second)
	if err != nil {
		t.Fatalf("Create Pool failled  = %v", err)
	}
	p.Init()
	p.TriggerError = resultOfError
	err = p.Add("test-1", test{"test-1", 200}, 1)

	time.Sleep(1 * time.Second)
}
func TestTriggerOk(t *testing.T) {
	var p, err = NewPool(test{}, "TEST", resultOfTest, 20, false, 1*time.Second)
	if err != nil {
		t.Fatalf("Create Pool failled  = %v", err)
	}
	p.Init()
	p.TriggerFinish = resultOnSuccess
	err = p.Add("test-1", test{"test-1", 100}, 1)

	time.Sleep(1 * time.Second)
}

func resultOfError(p *PoolEntry) error {
	fmt.Println("An error on ")
	return nil
}
func resultOnSuccess(p *PoolEntry) error {
	fmt.Println("Ok no error ")
	return nil
}
