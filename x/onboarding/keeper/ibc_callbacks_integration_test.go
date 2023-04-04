package keeper_test

import (
	"github.com/Canto-Network/Canto/v6/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/slices"
)

func FindEvent(events []sdk.Event, name string) sdk.Event {
	index := slices.IndexFunc(events, func(e sdk.Event) bool { return e.Type == name })
	if index == -1 {
		return sdk.Event{}
	}
	return events[index]
}

func ExtractAttributes(event sdk.Event) map[string]string {
	attrs := make(map[string]string)
	if event.Attributes == nil {
		return attrs
	}
	for _, a := range event.Attributes {
		attrs[string(a.Key)] = string(a.Value)
	}
	return attrs
}

var _ = Describe("Onboarding: Performing an IBC Transfer followed by autoswap and convert", Ordered, func() {
	coincanto := sdk.NewCoin("acanto", sdk.ZeroInt())
	coinUsdc := sdk.NewCoin("uUSDC", sdk.NewIntWithDecimal(10000, 6))
	coinAtom := sdk.NewCoin("uatom", sdk.NewIntWithDecimal(10000, 6))

	var (
		sender, receiver string
		senderAcc        sdk.AccAddress
		receiverAcc      sdk.AccAddress
		result           *sdk.Result
	)

	BeforeEach(func() {
		s.SetupTest()

	})

	Describe("from a non-authorized channel: Cosmos ---(uatom)---> Canto", func() {
		BeforeEach(func() {
			// send coins from Cosmos to canto
			sender = s.IBCCosmosChain.SenderAccount.GetAddress().String()
			receiver = s.cantoChain.SenderAccount.GetAddress().String()
			senderAcc = sdk.MustAccAddressFromBech32(sender)
			receiverAcc = sdk.MustAccAddressFromBech32(receiver)

			result = s.SendAndReceiveMessage(s.pathCosmoscanto, s.IBCCosmosChain, "uatom", 10000000000, sender, receiver, 1)

		})
		It("No swap and convert operation - acanto balance should be 0", func() {
			nativecanto := s.cantoChain.App.(*app.Canto).BankKeeper.GetBalance(s.cantoChain.GetContext(), receiverAcc, "acanto")
			Expect(nativecanto).To(Equal(coincanto))
		})
		It("Canto chain's IBC voucher balance should be same with the transferred amount", func() {
			ibcAtom := s.cantoChain.App.(*app.Canto).BankKeeper.GetBalance(s.cantoChain.GetContext(), receiverAcc, uatomIbcdenom)
			Expect(ibcAtom).To(Equal(sdk.NewCoin(uatomIbcdenom, coinAtom.Amount)))
			//atom := s.IBCCosmosChain.GetSimApp().BankKeeper.GetBalance(s.IBCCosmosChain.GetContext(), senderAcc, "uatom")
			//Expect(atom).To(Equal(sdk.NewCoin("uatom", sdk.ZeroInt())))
		})
		It("Cosmos chain's uatom balance should be 0", func() {
			atom := s.IBCCosmosChain.GetSimApp().BankKeeper.GetBalance(s.IBCCosmosChain.GetContext(), senderAcc, "uatom")
			Expect(atom).To(Equal(sdk.NewCoin("uatom", sdk.ZeroInt())))
		})
	})

	Describe("from an authorized channel: Gravity ---(uUSDC)---> Canto", func() {
		BeforeEach(func() {
			sender = s.IBCGravityChain.SenderAccount.GetAddress().String()
			receiver = s.cantoChain.SenderAccount.GetAddress().String()
			senderAcc = sdk.MustAccAddressFromBech32(sender)
			receiverAcc = sdk.MustAccAddressFromBech32(receiver)
		})

		Context("when no swap pool exists", func() {
			BeforeEach(func() {
				s.SendAndReceiveMessage(s.pathGravitycanto, s.IBCGravityChain, "uUSDC", 10000000000, sender, receiver, 1)
			})
			It("No swap and convert operation - acanto balance should be 0", func() {
				nativecanto := s.cantoChain.App.(*app.Canto).BankKeeper.GetBalance(s.cantoChain.GetContext(), receiverAcc, "acanto")
				Expect(nativecanto).To(Equal(coincanto))
			})
			It("Canto chain's IBC voucher balance should be same with the transferred amount", func() {
				ibcUsdc := s.cantoChain.App.(*app.Canto).BankKeeper.GetBalance(s.cantoChain.GetContext(), receiverAcc, uusdcIbcdenom)
				Expect(ibcUsdc).To(Equal(sdk.NewCoin(uusdcIbcdenom, coinUsdc.Amount)))
			})
		})

		Context("when swap pool exists", func() {
			BeforeEach(func() {
				s.CreatePool(uusdcIbcdenom)
			})
			When("ERC20 contract is not deployed", func() {
				When("acanto balance is 0 and not enough IBC token transferred to swap acanto", func() {
					BeforeEach(func() {
						result = s.SendAndReceiveMessage(s.pathGravitycanto, s.IBCGravityChain, "uUSDC", 1000000, sender, receiver, 1)
					})
					It("Balance of acanto should be same with the original acanto balance (0)", func() {
						nativecanto := s.cantoChain.App.(*app.Canto).BankKeeper.GetBalance(s.cantoChain.GetContext(), receiverAcc, "acanto")
						Expect(nativecanto).To(Equal(sdk.NewCoin("acanto", sdk.ZeroInt())))
					})
					It("Canto chain's IBC voucher balance should be same with the transferred amount", func() {
						ibcUsdc := s.cantoChain.App.(*app.Canto).BankKeeper.GetBalance(s.cantoChain.GetContext(), receiverAcc, uusdcIbcdenom)
						Expect(ibcUsdc).To(Equal(sdk.NewCoin(uusdcIbcdenom, sdk.NewInt(1000000))))
					})
					It("No ERC20 token pair exists", func() {
						pairID := s.cantoChain.App.(*app.Canto).Erc20Keeper.GetTokenPairID(s.cantoChain.GetContext(), uusdcIbcdenom)
						Expect(len(pairID)).To(Equal(0))
					})
				})

				When("Canto chain's acanto balance is 0", func() {
					BeforeEach(func() {
						result = s.SendAndReceiveMessage(s.pathGravitycanto, s.IBCGravityChain, "uUSDC", 10000000000, sender, receiver, 1)
					})
					It("Auto swap operation: balance of acanto should be same with the auto swap threshold", func() {
						autoSwapThreshold := s.cantoChain.App.(*app.Canto).OnboardingKeeper.GetParams(s.cantoChain.GetContext()).AutoSwapThreshold
						nativecanto := s.cantoChain.App.(*app.Canto).BankKeeper.GetBalance(s.cantoChain.GetContext(), receiverAcc, "acanto")
						Expect(nativecanto).To(Equal(sdk.NewCoin("acanto", autoSwapThreshold)))
					})
					It("Canto chain's IBC voucher balance should be less than the transferred amount (diff: swap amount)", func() {
						events := result.GetEvents()
						attrs := ExtractAttributes(FindEvent(events, "swap"))
						swapAmount, _ := sdk.NewIntFromString(attrs["amount"])

						ibcUsdc := s.cantoChain.App.(*app.Canto).BankKeeper.GetBalance(s.cantoChain.GetContext(), receiverAcc, uusdcIbcdenom)
						Expect(ibcUsdc).To(Equal(sdk.NewCoin(uusdcIbcdenom, coinUsdc.Amount.Sub(swapAmount))))
					})
					It("No ERC20 token pair exists", func() {
						pairID := s.cantoChain.App.(*app.Canto).Erc20Keeper.GetTokenPairID(s.cantoChain.GetContext(), uusdcIbcdenom)
						Expect(len(pairID)).To(Equal(0))
					})
				})

				When("Canto chain's acanto balance is between 0 and auto swap threshold (3canto)", func() {
					BeforeEach(func() {
						s.FundCantoChain(sdk.NewCoins(sdk.NewCoin("acanto", sdk.NewIntWithDecimal(3, 18))))
						result = s.SendAndReceiveMessage(s.pathGravitycanto, s.IBCGravityChain, "uUSDC", 10000000000, sender, receiver, 1)
					})
					It("Auto swap operation: balance of acanto should be same with the sum of original acanto balance and auto swap threshold", func() {
						autoSwapThreshold := s.cantoChain.App.(*app.Canto).OnboardingKeeper.GetParams(s.cantoChain.GetContext()).AutoSwapThreshold
						nativecanto := s.cantoChain.App.(*app.Canto).BankKeeper.GetBalance(s.cantoChain.GetContext(), receiverAcc, "acanto")
						Expect(nativecanto).To(Equal(sdk.NewCoin("acanto", autoSwapThreshold.Add(sdk.NewIntWithDecimal(3, 18)))))
					})
					It("Canto chain's IBC voucher balance should be less than the transferred amount (diff: swap amount)", func() {
						events := result.GetEvents()
						attrs := ExtractAttributes(FindEvent(events, "swap"))
						swapAmount, _ := sdk.NewIntFromString(attrs["amount"])

						ibcUsdc := s.cantoChain.App.(*app.Canto).BankKeeper.GetBalance(s.cantoChain.GetContext(), receiverAcc, uusdcIbcdenom)
						Expect(ibcUsdc).To(Equal(sdk.NewCoin(uusdcIbcdenom, coinUsdc.Amount.Sub(swapAmount))))
					})
					It("No ERC20 token pair exists", func() {
						pairID := s.cantoChain.App.(*app.Canto).Erc20Keeper.GetTokenPairID(s.cantoChain.GetContext(), uusdcIbcdenom)
						Expect(len(pairID)).To(Equal(0))
					})
				})

			})
			When("ERC20 contract is deployed", func() {
				BeforeEach(func() {
					s.setupRegisterCoin(metadataIbc)
				})
				When("Canto chain's acanto balance is 0", func() {
					BeforeEach(func() {
						result = s.SendAndReceiveMessage(s.pathGravitycanto, s.IBCGravityChain, "uUSDC", 10000000000, sender, receiver, 1)
					})
					It("Auto swap operation: balance of acanto should be same with the auto swap threshold", func() {
						autoSwapThreshold := s.cantoChain.App.(*app.Canto).OnboardingKeeper.GetParams(s.cantoChain.GetContext()).AutoSwapThreshold
						nativecanto := s.cantoChain.App.(*app.Canto).BankKeeper.GetBalance(s.cantoChain.GetContext(), receiverAcc, "acanto")
						Expect(nativecanto).To(Equal(sdk.NewCoin("acanto", autoSwapThreshold)))
					})
					It("Canto chain's IBC voucher balance should be 0 (all IBC voucher should be converted)", func() {
						ibcUsdc := s.cantoChain.App.(*app.Canto).BankKeeper.GetBalance(s.cantoChain.GetContext(), receiverAcc, uusdcIbcdenom)
						Expect(ibcUsdc).To(Equal(sdk.NewCoin(uusdcIbcdenom, sdk.ZeroInt())))
					})
					It("No ERC20 token pair exists", func() {
						pairID := s.cantoChain.App.(*app.Canto).Erc20Keeper.GetTokenPairID(s.cantoChain.GetContext(), uusdcIbcdenom)
						Expect(len(pairID)).To(Equal(0))
					})
				})
			})
		})
	})

})
