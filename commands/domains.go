package commands

import (
	"errors"
	"io"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
)

// Domain creates the domain commands heirarchy.
func Domain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "domain commands",
		Long:  "domain is used to access domain commands",
	}

	cmdDomainCreate := cmdBuilder(RunDomainCreate, "create", "create domain", writer, "c")
	cmd.AddCommand(cmdDomainCreate)
	addStringFlag(cmdDomainCreate, doctl.ArgDomainName, "", "Domain name")
	addStringFlag(cmdDomainCreate, doctl.ArgIPAddress, "", "IP address")

	cmdDomainList := cmdBuilder(RunDomainList, "list", "list comains", writer, "ls")
	cmd.AddCommand(cmdDomainList)

	cmdDomainGet := cmdBuilder(RunDomainGet, "get", "get domain", writer, "g")
	cmd.AddCommand(cmdDomainGet)
	addStringFlag(cmdDomainGet, doctl.ArgDomainName, "", "Domain name")

	cmdDomainDelete := cmdBuilder(RunDomainDelete, "delete", "delete droplet", writer, "g")
	cmd.AddCommand(cmdDomainDelete)
	addStringFlag(cmdDomainDelete, doctl.ArgDomainName, "", "Domain name")

	cmdRecord := &cobra.Command{
		Use:   "records",
		Short: "domain record commands",
		Long:  "commands for interacting with an individual domain",
	}
	cmd.AddCommand(cmdRecord)

	cmdRecordList := cmdBuilder(RunRecordList, "list", "list records", writer, "ls")
	cmdRecord.AddCommand(cmdRecordList)
	addStringFlag(cmdRecordList, doctl.ArgDomainName, "", "Domain name")

	cmdRecordCreate := cmdBuilder(RunRecordCreate, "create", "create record", writer, "c")
	cmdRecord.AddCommand(cmdRecordCreate)
	addStringFlag(cmdRecordCreate, doctl.ArgDomainName, "", "Domain name")
	addStringFlag(cmdRecordCreate, doctl.ArgRecordType, "", "Record type")
	addStringFlag(cmdRecordCreate, doctl.ArgRecordName, "", "Record name")
	addStringFlag(cmdRecordCreate, doctl.ArgRecordData, "", "Record data")
	addIntFlag(cmdRecordCreate, doctl.ArgRecordPriority, 0, "Record priority")
	addIntFlag(cmdRecordCreate, doctl.ArgRecordPort, 0, "Record port")
	addIntFlag(cmdRecordCreate, doctl.ArgRecordWeight, 0, "Record weight")

	cmdRecordDelete := cmdBuilder(RunRecordDelete, "delete", "delete record", writer, "d")
	cmdRecord.AddCommand(cmdRecordDelete)
	addStringFlag(cmdRecordDelete, doctl.ArgDomainName, "", "Domain name")
	addIntFlag(cmdRecordDelete, doctl.ArgRecordID, 0, "Record ID")

	cmdRecordUpdate := cmdBuilder(RunRecordUpdate, "update", "update record", writer, "u")
	cmdRecord.AddCommand(cmdRecordUpdate)
	addStringFlag(cmdRecordUpdate, doctl.ArgDomainName, "", "Domain name")
	addIntFlag(cmdRecordUpdate, doctl.ArgRecordID, 0, "Record ID")
	addStringFlag(cmdRecordUpdate, doctl.ArgRecordType, "", "Record type")
	addStringFlag(cmdRecordUpdate, doctl.ArgRecordName, "", "Record name")
	addStringFlag(cmdRecordUpdate, doctl.ArgRecordData, "", "Record data")
	addIntFlag(cmdRecordUpdate, doctl.ArgRecordPriority, 0, "Record priority")
	addIntFlag(cmdRecordUpdate, doctl.ArgRecordPort, 0, "Record port")
	addIntFlag(cmdRecordUpdate, doctl.ArgRecordWeight, 0, "Record weight")

	return cmd
}

// RunDomainCreate runs domain create.
func RunDomainCreate(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	req := &godo.DomainCreateRequest{
		Name:      doctl.DoctlConfig.GetString(ns, "domain-name"),
		IPAddress: doctl.DoctlConfig.GetString(ns, "ip-address"),
	}

	d, _, err := client.Domains.Create(req)
	if err != nil {
		return err
	}

	return doctl.DisplayOutput(d, out)
}

// RunDomainList runs domain create.
func RunDomainList(cmdName string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()

	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := client.Domains.List(opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := doctl.PaginateResp(f)
	if err != nil {
		return err
	}

	list := make([]godo.Domain, len(si))
	for i := range si {
		list[i] = si[i].(godo.Domain)
	}

	return doctl.DisplayOutput(list, out)
}

// RunDomainGet retrieves a domain by name.
func RunDomainGet(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	id := doctl.DoctlConfig.GetString(ns, doctl.ArgDomainName)

	if len(id) < 1 {
		return errors.New("invalid domain name")
	}

	d, _, err := client.Domains.Get(id)
	if err != nil {
		return err
	}

	return doctl.DisplayOutput(d, out)
}

// RunDomainDelete deletes a domain by name.
func RunDomainDelete(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	name := doctl.DoctlConfig.GetString(ns, doctl.ArgDomainName)

	if len(name) < 1 {
		return errors.New("invalid domain name")
	}

	_, err := client.Domains.Delete(name)
	return err
}

// RunRecordList list records for a domain.
func RunRecordList(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	name := doctl.DoctlConfig.GetString(ns, doctl.ArgDomainName)

	if len(name) < 1 {
		return errors.New("domain name is missing")
	}

	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := client.Domains.Records(name, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := doctl.PaginateResp(f)
	if err != nil {
		return err
	}

	list := make([]godo.DomainRecord, len(si))
	for i := range si {
		list[i] = si[i].(godo.DomainRecord)
	}

	return doctl.DisplayOutput(list, out)
}

// RunRecordCreate creates a domain record.
func RunRecordCreate(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	name := doctl.DoctlConfig.GetString(ns, doctl.ArgDomainName)

	drcr := &godo.DomainRecordEditRequest{
		Type:     doctl.DoctlConfig.GetString(ns, doctl.ArgRecordType),
		Name:     doctl.DoctlConfig.GetString(ns, doctl.ArgRecordName),
		Data:     doctl.DoctlConfig.GetString(ns, doctl.ArgRecordData),
		Priority: doctl.DoctlConfig.GetInt(ns, doctl.ArgRecordPriority),
		Port:     doctl.DoctlConfig.GetInt(ns, doctl.ArgRecordPort),
		Weight:   doctl.DoctlConfig.GetInt(ns, doctl.ArgRecordWeight),
	}

	if len(drcr.Type) == 0 {
		return errors.New("record request is missing type")
	}

	r, _, err := client.Domains.CreateRecord(name, drcr)
	if err != nil {
		return err
	}

	return doctl.DisplayOutput(r, out)
}

// RunRecordDelete deletes a domain record.
func RunRecordDelete(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	domainName := doctl.DoctlConfig.GetString(ns, doctl.ArgDomainName)
	recordID := doctl.DoctlConfig.GetInt(ns, doctl.ArgRecordID)

	_, err := client.Domains.DeleteRecord(domainName, recordID)
	return err
}

// RunRecordUpdate updates a domain record.
func RunRecordUpdate(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	domainName := doctl.DoctlConfig.GetString(ns, doctl.ArgDomainName)
	recordID := doctl.DoctlConfig.GetInt(ns, doctl.ArgRecordID)

	drcr := &godo.DomainRecordEditRequest{
		Type:     doctl.DoctlConfig.GetString(ns, doctl.ArgRecordType),
		Name:     doctl.DoctlConfig.GetString(ns, doctl.ArgRecordName),
		Data:     doctl.DoctlConfig.GetString(ns, doctl.ArgRecordData),
		Priority: doctl.DoctlConfig.GetInt(ns, doctl.ArgRecordPriority),
		Port:     doctl.DoctlConfig.GetInt(ns, doctl.ArgRecordPort),
		Weight:   doctl.DoctlConfig.GetInt(ns, doctl.ArgRecordWeight),
	}

	r, _, err := client.Domains.EditRecord(domainName, recordID, drcr)
	if err != nil {
		return err
	}

	return doctl.DisplayOutput(r, out)
}
