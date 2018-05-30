package accounts

// Accounts
// Account
// Copyright © 2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
)

/*
newAccount - create new Account.
*/
func New() *Accounts {
	return &Accounts{
		data: make(map[string]*Account),
	}
}

type Accounts struct {
	data map[string]*Account
}

func (a *Accounts) Account(id string) *Account {
	acc, ok := a.data[id]
	if !ok {
		acc = newAccount()
	}
	return acc
}

type Balance struct {
	available uint64
	blocked   uint64
}

type Account struct {
	Balance
	//available uint64
	//blocked   uint64
	blocks map[string]uint64
}

/*
newAccount - create new Account.
*/
func newAccount() *Account {
	return &Account{
		blocks: make(map[string]uint64),
	}
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

func (a *Account) BlockNoFix(amount uint64) (uint64, uint64, error) {
	if a.available < amount {
		return a.available, a.blocked, fmt.Errorf("Blocking error - there is %d, but blocked %d.", a.available, amount)
	}
	newAviable := a.available - amount
	newBlocked := a.blocked + amount
	if newBlocked < a.blocked {
		return a.available, a.blocked, fmt.Errorf("Overflow error: there is %d, add %d, get %d. (Block operation)", a.blocked, amount, newBlocked)
	}
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

func (a *Account) UnblockNoFix(amount uint64) (uint64, uint64, error) {
	newAviable := a.available + amount
	newBlocked := a.blocked - amount
	if newAviable < a.available {
		return a.available, a.blocked, fmt.Errorf("Overflow error: there is %d, add %d, get %d. (Unlock operation)", a.available, amount, newAviable)
	}
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
