package accounter

// Accounts
// Account
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
)

type Account1 struct {
	data map[string]int64
}

func (a *Account1) AddAmount(sub string, amount int64) int64 {
	cur, _ := a.data[sub]
	//if !ok {
	//	cur = 0
	//}
	cur += amount
	if cur > 0 {
		a.data[sub] = cur
	}
	return cur
}

type Account struct {
	available uint64
	blocked   uint64
	blocks    map[string]uint64
}

/*
NewAccount - create new Account.
*/
func NewAccount() *Account {
	return &Account{}
}

func (a *Account) Debit(amount uint64) (uint64, uint64, error) {
	newAviable := a.available + amount
	if newAviable < a.available {
		return a.available, a.blocked, fmt.Errorf("Overflow error: there is %d, add %d, get %d. (Debit operation)", a.available, amount, newAviable)
	}
	a.available = newAviable
	return a.available, a.blocked, nil
}

func (a *Account) Block(key string, amount uint64) (uint64, uint64, error) {
	if _, ok := a.blocks[key]; ok {
		return a.available, a.blocked, fmt.Errorf("This key is already taken.")
	}
	if a.available < amount {
		return a.available, a.blocked, fmt.Errorf("Blocking error - there is %d, but blocked %d.", a.available, amount)
	}

	newAviable := a.available - amount
	newBlocked := a.blocked + amount
	if newBlocked < a.blocked {
		return a.available, a.blocked, fmt.Errorf("Overflow error: there is %d, add %d, get %d. (Block operation)", a.blocked, amount, newBlocked)
	}
	a.blocks[key] = amount
	a.available = newAviable
	a.blocked = newBlocked
	return a.available, a.blocked, nil
}

func (a *Account) Unblock(key string, amount uint64) (uint64, uint64, error) {
	sum, ok := a.blocks[key]
	if !ok {
		return a.available, a.blocked, fmt.Errorf("This key is missing.")
	}
	if sum != amount {
		return a.available, a.blocked, fmt.Errorf("The amount does not match the blocked amount..")
	}
	newAviable := a.available + amount
	newBlocked := a.blocked - amount

	if newAviable < a.available {
		return a.available, a.blocked, fmt.Errorf("Overflow error: there is %d, add %d, get %d. (Unlock operation)", a.available, amount, newAviable)
	}
	delete(a.blocks, key)
	a.available = newAviable
	a.blocked = newBlocked
	return a.available, a.blocked, nil
}

func (a *Account) Credit(key string, amount uint64) (uint64, uint64, error) {
	sum, ok := a.blocks[key]
	if !ok {
		return a.available, a.blocked, fmt.Errorf("This key is missing.")
	}
	if sum != amount {
		return a.available, a.blocked, fmt.Errorf("The amount does not match the blocked amount..")
	}
	newBlocked := a.blocked - amount

	delete(a.blocks, key)
	a.blocked = newBlocked
	return a.available, a.blocked, nil
}

func (a *Account) WriteOff(amount uint64) (uint64, uint64, error) { //  Credit operation without intermediate blocking of funds.
	if a.available < amount {
		return a.available, a.blocked, fmt.Errorf("Blocking error - there is %d, but blocked %d.", a.available, amount)
	}
	newAviable := a.available - amount
	if newAviable > a.available {
		return a.available, a.blocked, fmt.Errorf("Overflow error: there is %d, add %d, get %d. (WriteOff operation)", a.available, amount, newAviable)
	}
	a.available = newAviable
	return a.available, a.blocked, nil
}
