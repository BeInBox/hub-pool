package POOL

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type test struct {
	Name string
	Val  int
}

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
func resultOfTest(p *poolEntry) error {
	nb := rand.Intn(200)
	fmt.Println("JE demarre", p.Key, "avec une attente de ", nb)

	time.Sleep(time.Duration(nb) * time.Millisecond)
	if nb > 170 {
		fmt.Println("Merde ca marcher pas")
		return errors.New("Ko")
	}
	fmt.Println("Ok ca marcher")
	return nil
}
