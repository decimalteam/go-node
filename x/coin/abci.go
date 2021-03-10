package coin

import (
	"bitbucket.org/decimalteam/go-node/utils/updates"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	k.ClearCoinCache()

	if ctx.BlockHeight() == updates.Update7Block {
		coins := map[string]string{
			"bloggcoin":  "dx108sl80sdra0etw3fxsrd0euvdcer4j89xxnjzm",
			"btt":        "dx16rr3cvdgj8jsywhx8lfteunn9uz0xg2c7ua9nl",
			"convenienc": "dx1d4gd9vgu25yl8jegh0c5z06h95vs4apcs430zr",
			"danielcoin": "dx1u2st9ucd8m67tkqg8a599xh4kleszzzv5z6qjv",
			"dynamics":   "dx1d4gd9vgu25yl8jegh0c5z06h95vs4apcs430zr",
			"ethereum":   "dx1py8nz9e5j2ppake57rfu7e45dcvm6vs7gg79tg",
			"finaltest":  "dx1l3mlmzl0jwz05pkwlkwana6ruqupryeztr5q4f",
			"fivegrant":  "dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v",
			"fourgrant":  "dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v",
			"gfhgfhgfhg": "dx1hv8sx65303fyq4hp2h6j0rndkkqfxghvxz43jw",
			"investcoin": "dx1hv8sx65303fyq4hp2h6j0rndkkqfxghvxz43jw",
			"legiondocs": "dx1566ntqfrh9mzt9ldqyarfrq5ljkxjes897cp7c",
			"legiong2":   "dx1k2ppgdql2slvl06g926pm8dfszh69wywlwzdww",
			"legiongame": "dx1m4ta03kx3m3yd3ttm7hcpg6d54h2vy06axsd0a",
			"legionsert": "dx1l3mlmzl0jwz05pkwlkwana6ruqupryeztr5q4f",
			"legiontest": "dx1k2ppgdql2slvl06g926pm8dfszh69wywlwzdww",
			"nazvala":    "dx1jxyrcczw2e4dhrlf52e6mlgap5q99vjqsjjyun",
			"nazvaniemo": "dx1jxyrcczw2e4dhrlf52e6mlgap5q99vjqsjjyun",
			"nedviga":    "dx108sl80sdra0etw3fxsrd0euvdcer4j89xxnjzm",
			"onegrant":   "dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v",
			"privetuli":  "dx1jxyrcczw2e4dhrlf52e6mlgap5q99vjqsjjyun",
			"proverkamo": "dx1h4a95lgvskagn8grgnx58erk5sl3ewy8ktyk5h",
			"pumpcoin25": "dx1mn544x23s64pxq686pg7vkrkrgunypztul263k",
			"pumpsuperc": "dx1mn544x23s64pxq686pg7vkrkrgunypztul263k",
			"qwertyiopo": "dx1gjs3rpcelrdg6yhrasdg9cgv4ms8l697jzklpv",
			"renatadimo": "dx1hz2y5rch7v39jk9l0suzg3ky8uhct9r922rp5p",
			"secondcoin": "dx1mvqrrrlcd0gdt256jxg7n68e4neppu5t24e8h6",
			"sertificar": "dx1k2ppgdql2slvl06g926pm8dfszh69wywlwzdww",
			"sertificat": "dx1hv8sx65303fyq4hp2h6j0rndkkqfxghvxz43jw",
			"sertifikat": "dx1k2ppgdql2slvl06g926pm8dfszh69wywlwzdww",
			"sfsfsfsf":   "dx1cfvvmrg5avgdx2gpr8ffddgceev74zz4rt52jl",
			"superpump":  "dx1mn544x23s64pxq686pg7vkrkrgunypztul263k",
			"superstar1": "dx1hz2y5rch7v39jk9l0suzg3ky8uhct9r922rp5p",
			"tengrants":  "dx1nrr6er27mmcufmaqm4dyu6c5r6489cfm35m4ft",
			"testtt":     "dx1gtlgwrnads2xh7uydlg6pa5htjmqgf69xjfgcf",
			"thepercent": "dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v",
			"threegrant": "dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v",
			"twogrant":   "dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v",
			"youblogger": "dx1hz2y5rch7v39jk9l0suzg3ky8uhct9r922rp5p",
			"zxczxc01":   "dx185xupdr6w4klhu5js4j7zdwswcun0m4k27npvg",
		}
		for symbol, address := range coins {
			coin, err := k.GetCoin(ctx, symbol)
			if err != nil {
				panic(err)
			}
			coin.Creator, err = sdk.AccAddressFromBech32(address)
			if err != nil {
				panic(err)
			}
			k.SetCoin(ctx, coin)
		}
	}
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k Keeper) {
}
