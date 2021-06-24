package cmd

import (
	"errors"
	"fmt"
	"github.com/roseboy/httpcase/internal/httpcase"
	"github.com/roseboy/httpcase/util"
	"github.com/spf13/cobra"
)

type runCmd struct {
	cmd  *cobra.Command
	opts runOpts
}

type runOpts struct {
	env string
	out string
}

func newRunCmd() *runCmd {
	root := &runCmd{}

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run api test case file",
		Run: func(cmd *cobra.Command, args []string) {
			err := doThat(root, args)
			if err != nil {
				util.Println(util.Red("Error:"), err.Error())
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a path argument")
			}
			if !util.IsExist(args[0]) {
				return fmt.Errorf("invalid file path \"%s\"", args[0])
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&root.opts.env, "env", "e", "", "env flag")
	cmd.Flags().StringVarP(&root.opts.out, "out", "o", "", "output test report file")

	root.cmd = cmd
	return root
}

func doThat(cmd *runCmd, args []string) error {
	var (
		path = args[0]
		tags string
	)

	if len(args) > 1 {
		tags = args[1]
	}

	testCtx := httpcase.NewTestContext()
	testCtx.Env = cmd.opts.env
	testCtx.Out = cmd.opts.out
	testCtx.Tags = tags

	//读取文件
	codes, err := httpcase.ReadCaseFile(path)
	if err != nil {
		return err
	}

	// 解析测试文件
	testRequest, err := httpcase.NewCompiler(testCtx, codes).Compile()
	if err != nil {
		return err
	}

	// 运行测试
	testResult, err := testCtx.Init(testRequest).Run()
	if err != nil {
		return err
	}

	// 输出报告
	err = httpcase.PrintReport(testCtx, testResult)
	if err != nil {
		return err
	}

	return nil
}
