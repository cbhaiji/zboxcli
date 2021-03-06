package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

func printStakePoolInfo(info *sdk.StakePoolInfo) {
	fmt.Println("pool_id:", info.ID)
	fmt.Println("balance:", info.Balance)

	fmt.Println("capacity:")
	fmt.Println("  free:       ", info.Free, "(for current write price)")
	fmt.Println("  capacity:   ", info.Capacity, "(blobber bid)")
	fmt.Println("  write_price:", info.WritePrice, "(blobber write price)")

	if len(info.Offers) == 0 {
		fmt.Println("offers: no opened offers")
	} else {
		fmt.Println("offers:")
		for _, off := range info.Offers {
			fmt.Println("- lock:      ", off.Lock)
			fmt.Println("  expire:    ", off.Expire.ToTime())
			fmt.Println("  allocation:", off.AllocationID)
			fmt.Println("  expired:   ", off.IsExpired)
		}
		fmt.Println("offers_total:", info.OffersTotal, "(held by opened offers)")
	}

	if len(info.Delegate) == 0 {
		fmt.Println("delegate_pools: no delegate pools")
	} else {
		fmt.Println("delegate_pools:")
		for _, dp := range info.Delegate {
			fmt.Println("- id:         ", dp.ID)
			fmt.Println("  balance:    ", dp.Balance)
			fmt.Println("  delegate_id:", dp.DelegateID)
			fmt.Println("  earnings:   ", dp.Earnings, "(payed interests for the delegate pool)")
			fmt.Println("  penalty:    ", dp.Penalty, "(penalty for the delegate pool)")
			fmt.Println("  interests:  ", dp.Interests, "(interests not payed yet, can be given by 'sp-pay-interests' command)")
		}
	}
	fmt.Println("earnings:", info.Earnings, "(total interests earnings for all delegate pools for all time)")
	fmt.Println("penalty:", info.Penalty, "(total blobber penalty for all time)")

	fmt.Println("rewards: (excluding interests)")
	fmt.Println("  balance:  ", info.Rewards.Balance, "(current rewards can be unlocked)")
	fmt.Println("  blobber:  ", info.Rewards.Blobber, "(for all time)")
	fmt.Println("  validator:", info.Rewards.Validator, "(for all time)")

}

// spInfo information
var spInfo = &cobra.Command{
	Use:   "sp-info",
	Short: "Stake pool information.",
	Long:  `Stake pool information.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags     = cmd.Flags()
			blobberID string
			err       error
		)

		if flags.Changed("blobber_id") {
			if blobberID, err = flags.GetString("blobber_id"); err != nil {
				log.Fatalf("can't get 'blobber_id' flag: %v", err)
			}
		}

		var info *sdk.StakePoolInfo
		if info, err = sdk.GetStakePoolInfo(blobberID); err != nil {
			log.Fatalf("Failed to get stake pool info: %v", err)
		}
		printStakePoolInfo(info)
	},
}

// spLock locks tokens a stake pool lack
var spLock = &cobra.Command{
	Use:   "sp-lock",
	Short: "Lock tokens lacking in stake pool.",
	Long:  `Lock tokens lacking in stake pool.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags     = cmd.Flags()
			blobberID string
			tokens    float64
			fee       float64
			err       error
		)

		if flags.Changed("blobber_id") {
			if blobberID, err = flags.GetString("blobber_id"); err != nil {
				log.Fatalf("invalid 'blobber_id' flag: %v", err)
			}
		}

		if !flags.Changed("tokens") {
			log.Fatal("missing required 'tokens' flag")
		}

		if tokens, err = flags.GetFloat64("tokens"); err != nil {
			log.Fatal("invalid 'tokens' flag: ", err)
		}

		if flags.Changed("fee") {
			if fee, err = flags.GetFloat64("fee"); err != nil {
				log.Fatal("invalid 'fee' flag: ", err)
			}
		}

		var poolID string
		poolID, err = sdk.StakePoolLock(blobberID,
			zcncore.ConvertToValue(tokens), zcncore.ConvertToValue(fee))
		if err != nil {
			log.Fatalf("Failed to lock tokens in stake pool: %v", err)
		}
		fmt.Println("tokens locked, pool id:", poolID)
	},
}

// spUnlock unlocks tokens in stake pool
var spUnlock = &cobra.Command{
	Use:   "sp-unlock",
	Short: "Unlock tokens in stake pool.",
	Long:  `Unlock tokens in stake pool.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags             = cmd.Flags()
			blobberID, poolID string
			fee               float64
			err               error
		)

		if flags.Changed("blobber_id") {
			if blobberID, err = flags.GetString("blobber_id"); err != nil {
				log.Fatalf("invalid 'blobber_id' flag: %v", err)
			}
		}

		if !flags.Changed("pool_id") {
			log.Fatal("missing required 'pool_id' flag")
		}

		if poolID, err = flags.GetString("pool_id"); err != nil {
			log.Fatal("invalid 'pool_id' flag: ", err)
		}

		if flags.Changed("fee") {
			if fee, err = flags.GetFloat64("fee"); err != nil {
				log.Fatal("invalid 'fee' flag: ", err)
			}
		}

		err = sdk.StakePoolUnlock(blobberID, poolID,
			zcncore.ConvertToValue(fee))
		if err != nil {
			log.Fatalf("Failed to unlock tokens in stake pool: %v", err)
		}
		fmt.Println("tokens has unlocked, pool deleted")
	},
}

// spTakeRewards unlocks rewards of the blobber including
// validator rewards, and excluding interests
var spTakeRewards = &cobra.Command{
	Use:   "sp-take-rewards",
	Short: "Take blobber rewards.",
	Long: `Take blobber rewards, including all blobber rewards, rewards of
related validator, and excluding interests.`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags     = cmd.Flags()
			blobberID string
			err       error
		)

		if flags.Changed("blobber_id") {
			if blobberID, err = flags.GetString("blobber_id"); err != nil {
				log.Fatalf("invalid 'blobber_id' flag: %v", err)
			}
		}

		if err = sdk.StakePoolTakeRewards(blobberID); err != nil {
			log.Fatalf("Failed to take rewards: %v", err)
		}
		fmt.Println("rewards has taken")
	},
}

// spPayInterests pays interests not payed yet. A stake pool changes
// pays all interests can be payed. But if stake pool is not changed,
// then user can manually pay the interests.
var spPayInterests = &cobra.Command{
	Use:   "sp-pay-interests",
	Short: "Pay interests not payed yet.",
	Long:  `Pay interests not payed.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags     = cmd.Flags()
			blobberID string
			err       error
		)

		if flags.Changed("blobber_id") {
			if blobberID, err = flags.GetString("blobber_id"); err != nil {
				log.Fatalf("invalid 'blobber_id' flag: %v", err)
			}
		}

		if err = sdk.StakePoolPayInterests(blobberID); err != nil {
			log.Fatalf("Failed to pay interests: %v", err)
		}
		fmt.Println("interests has payed")
	},
}

func init() {
	rootCmd.AddCommand(spInfo)
	rootCmd.AddCommand(spLock)
	rootCmd.AddCommand(spUnlock)
	rootCmd.AddCommand(spTakeRewards)
	rootCmd.AddCommand(spPayInterests)

	spInfo.PersistentFlags().String("blobber_id", "",
		"for given blobber, default is current client")

	spLock.PersistentFlags().String("blobber_id", "",
		"for given blobber, default is current client")
	spLock.PersistentFlags().Float64("tokens", 0.0,
		"tokens to lock, required")
	spLock.PersistentFlags().Float64("fee", 0.0,
		"transaction fee, default 0")
	spLock.MarkFlagRequired("tokens")

	spUnlock.PersistentFlags().String("blobber_id", "",
		"for given blobber, default is current client")
	spUnlock.PersistentFlags().String("pool_id", "",
		"pool id to unlock")
	spUnlock.PersistentFlags().Float64("fee", 0.0,
		"transaction fee, default 0")
	spUnlock.MarkFlagRequired("tokens")
	spUnlock.MarkFlagRequired("pool_id")

	spTakeRewards.PersistentFlags().String("blobber_id", "",
		"for given blobber, default is current client")

	spPayInterests.PersistentFlags().String("blobber_id", "",
		"for given blobber, default is current client")
}
