package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/digitalocean/doctl"
	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// SSHKeys creates the ssh key commands heirarchy.
func SSHKeys() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sshkey",
		Aliases: []string{"k"},
		Short:   "sshkey commands",
		Long:    "sshkey is used to access ssh key commands",
	}

	cmdSSHKeysList := cmdBuilder(RunKeyList, "list", "list ssh keys", writer, "ls")
	cmd.AddCommand(cmdSSHKeysList)

	cmdSSHKeysGet := cmdBuilder(RunKeyGet, "get", "get ssh key", writer, "g")
	cmd.AddCommand(cmdSSHKeysGet)
	addStringFlag(cmdSSHKeysGet, doctl.ArgKey, "", "Key ID or fingerprint")

	cmdSSHKeysCreate := cmdBuilder(RunKeyCreate, "create", "create ssh key", writer, "c")
	cmd.AddCommand(cmdSSHKeysCreate)
	addStringFlag(cmdSSHKeysCreate, doctl.ArgKeyName, "", "Key name")
	addStringFlag(cmdSSHKeysCreate, doctl.ArgKeyPublicKey, "", "Key contents")

	cmdSSHKeysImport := cmdBuilder(RunKeyImport, "import", "import ssh key", writer, "i")
	cmd.AddCommand(cmdSSHKeysImport)
	addStringFlag(cmdSSHKeysImport, doctl.ArgKeyName, "", "Key name")
	addStringFlag(cmdSSHKeysImport, doctl.ArgKeyPublicKeyFile, "", "Public key file")

	cmdSSHKeysDelete := cmdBuilder(RunKeyDelete, "delete", "delete ssh key", writer, "d")
	cmd.AddCommand(cmdSSHKeysDelete)
	addStringFlag(cmdSSHKeysDelete, doctl.ArgKey, "", "Key ID or fingerprint")

	cmdSSHKeysUpdate := cmdBuilder(RunKeyUpdate, "update", "update ssh key", writer, "u")
	cmd.AddCommand(cmdSSHKeysUpdate)
	addStringFlag(cmdSSHKeysUpdate, doctl.ArgKey, "", "Key ID or fingerprint")
	addStringFlag(cmdSSHKeysUpdate, doctl.ArgKeyName, "", "Key name")

	return cmd
}

// RunKeyList lists keys.
func RunKeyList(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()

	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := client.Keys.List(opt)
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

	list := make([]godo.Key, len(si))
	for i := range si {
		list[i] = si[i].(godo.Key)
	}

	return doctl.DisplayOutput(list, out)
}

// RunKeyGet retrieves a key.
func RunKeyGet(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	rawKey := doctl.DoctlConfig.GetString(ns, doctl.ArgKey)

	var err error
	var key *godo.Key
	if i, aerr := strconv.Atoi(rawKey); aerr == nil {
		key, _, err = client.Keys.GetByID(i)
	} else {
		if len(rawKey) > 0 {
			key, _, err = client.Keys.GetByFingerprint(rawKey)
		} else {
			err = fmt.Errorf("missing key id or fingerprint")
		}
	}

	if err != nil {
		return err
	}

	return doctl.DisplayOutput(key, out)
}

// RunKeyCreate uploads a SSH key.
func RunKeyCreate(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()

	kcr := &godo.KeyCreateRequest{
		Name:      doctl.DoctlConfig.GetString(ns, doctl.ArgKeyName),
		PublicKey: doctl.DoctlConfig.GetString(ns, doctl.ArgKeyPublicKey),
	}

	r, _, err := client.Keys.Create(kcr)
	if err != nil {
		logrus.WithField("err", err).Fatal("could not create key")
	}

	return doctl.DisplayOutput(r, out)
}

// RunKeyImport imports a key from a file
func RunKeyImport(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()

	keyPath := doctl.DoctlConfig.GetString(ns, doctl.ArgKeyPublicKeyFile)
	keyName := doctl.DoctlConfig.GetString(ns, doctl.ArgKeyName)

	keyFile, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return err
	}

	_, comment, _, _, err := ssh.ParseAuthorizedKey(keyFile)
	if err != nil {
		return err
	}

	if len(keyName) < 1 {
		keyName = comment
	}

	kcr := &godo.KeyCreateRequest{
		Name:      keyName,
		PublicKey: string(keyFile),
	}

	r, _, err := client.Keys.Create(kcr)
	if err != nil {
		return err
	}

	return doctl.DisplayOutput(r, out)
}

// RunKeyDelete deletes a key.
func RunKeyDelete(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	rawKey := doctl.DoctlConfig.GetString(ns, doctl.ArgKey)

	var err error
	if i, aerr := strconv.Atoi(rawKey); aerr == nil {
		_, err = client.Keys.DeleteByID(i)
	} else {
		_, err = client.Keys.DeleteByFingerprint(rawKey)
	}

	return err
}

// RunKeyUpdate updates a key.
func RunKeyUpdate(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	rawKey := doctl.DoctlConfig.GetString(ns, doctl.ArgKey)

	req := &godo.KeyUpdateRequest{
		Name: doctl.DoctlConfig.GetString(ns, doctl.ArgKeyName),
	}

	var err error
	var key *godo.Key
	if i, aerr := strconv.Atoi(rawKey); aerr == nil {
		key, _, err = client.Keys.UpdateByID(i, req)
	} else {
		key, _, err = client.Keys.UpdateByFingerprint(rawKey, req)
	}

	if err != nil {
		return err
	}

	return doctl.DisplayOutput(key, out)
}
