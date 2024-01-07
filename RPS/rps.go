package RPS

import (

	"github.com/nspcc-dev/neo-go/pkg/interop"
	"github.com/nspcc-dev/neo-go/pkg/interop/contract"
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
	"github.com/nspcc-dev/neo-go/pkg/interop/storage"
)

const (
	gasDecimals    = 1_0000_0000
	initialBalance = 3000
	zaCoinHashKey = "zaCoinHash"
)

type Result struct {
	win bool
	tie bool
	lose bool
}

func _deploy(data interface{}, isUpdate bool) {
	if isUpdate {
		return
	}

	// Parse hash of forint contract from incoming data
	args := data.(struct {
		zaCoinHash interop.Hash160
	})

	if len(args.zaCoinHash) != interop.Hash160Len {
                panic("invalid hash of zaCoin contract")
        }


	ctx := storage.GetContext()
	storage.Put(ctx, zaCoinHashKey, args.zaCoinHash)
}

func PlayRPS(playerChoice string) {
	
	computerChoice := (runtime.GetRandom() % 3) + 1
	if computerChoice == 1 {
		computerChoice := "rock"
	} else if computerChoice == 2 {
		computerChoice := "paper"
	} else {
		computerChoice := "scissors"
	}

	result := isWinner(playerChoice, computerChoice)

    if result.tie {
		return
    } else if result.win {
        changePlayerBalance(ctx, playerOwner, bet)
    } else {
        changePlayerBalance(ctx, playerOwner, -bet)
    }
}

func isWinner(playerChoice, computerChoice string) Result {

	if playerChoice == computerChoice {
        return Result{tie: true}
    }

    if playerChoice == "rock" && computerChoice == "scissors" {
        return Result{win: true}
    }

    if playerChoice == "scissors" && computerChoice == "paper" {
        return Result{win: true}
    }

    if playerChoice == "paper" && computerChoice == "rock" {
        return Result{win: true}
    }

    return Result{lose: true}
}


func OnNEP17Payment(from interop.Hash160, amount int, data any) {
	ctx := storage.GetContext()
	zaCoinHash := storage.Get(ctx, zaCoinHashKey).(interop.Hash160)

	callingHash := runtime.GetCallingScriptHash()
	if !callingHash.Equals(zaCoinHash) {
		panic("only ZC is accepted")
	}
}

func changePlayerBalance(ctx storage.Context, playerOwner interop.Hash160, balanceChange int) {
	zaCoinHash := storage.Get(ctx, zaCoinHashKey).(interop.Hash160)
	playerContract := runtime.GetExecutingScriptHash()

	var from, to interop.Hash160
	var transferAmount int
	if balanceChange > 0 {
		// Transfer funds from contract to player owner
		from = playerContract
		to = playerOwner
		transferAmount = balanceChange
	} else {
		// Transfer funds from player owner to contract
		from = playerOwner
		to = playerContract
		transferAmount = -balanceChange // We flip sender/receiver, but keep amount positive
	}

	transferred := contract.Call(zaCoinHash, "transfer", contract.All, from, to, transferAmount, nil).(bool)
	if !transferred {
		panic("failed to transfer zaCoins")
	}
}