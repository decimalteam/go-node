package cli

import (
	"encoding/json"
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
)

func parseSubmitProposalFlags(fs *flag.FlagSet) (*proposal, error) {
	proposal := &proposal{}
	proposalFile, _ := fs.GetString(FlagProposal)

	if proposalFile == "" {
		proposal.Title = viper.GetString(FlagTitle)
		proposal.Description = viper.GetString(FlagDescription)
		proposal.VotingStartBlock = viper.GetUint64(FlagVotingStartBlock)
		proposal.VotingEndBlock = viper.GetUint64(FlagVotingEndBlock)
		return proposal, nil
	}

	for _, flag := range ProposalFlags {
		if viper.GetString(flag) != "" {
			return nil, fmt.Errorf("--%s flag provided alongside --proposal, which is a noop", flag)
		}
	}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, proposal)
	if err != nil {
		return nil, err
	}

	return proposal, nil
}
