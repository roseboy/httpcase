package cmd

import (
	"fmt"
	"github.com/roseboy/httpcase/internal/httpcase"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return doThat(root, args)
		},
	}

	cmd.Flags().StringVarP(&root.opts.env, "env", "e", "", "env flag")
	cmd.Flags().StringVarP(&root.opts.out, "out", "o", "", "output test report file")

	root.cmd = cmd
	return root
}

func doThat(cmd *runCmd, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("can't get test case file")
	}
	path := args[0]
	testCtx := httpcase.NewTestContext()
	testCtx.Env = cmd.opts.env
	testCtx.Out = cmd.opts.out

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
