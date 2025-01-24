package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/checker"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/manager"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/initialization"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/process"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/statistics"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/formatutil"
)

const runFailed = -1

func Init(ctx context.Context) error {
	initialization.InitFlag()
	flag.Parse()
	initialization.SetFlagLimit()
	initialization.GenUserIDs()
	initialization.InitGlobalData()
	if err := initialization.InitLog(config.GetConf()); err != nil {
		return err
	}
	initialization.PrintFlag()

	return nil
}

func DoFlagFunc(ctx context.Context) (err error) {
	defer func() {
		// capture check err
		ch := checker.CloseAndGetCheckErrChan()
		for e := range ch {
			if err == nil {
				err = e
			}
			log.ZError(ctx, "check err", err)
		}
	}()
	var (
		pro = process.NewProcess()

		m           = manager.NewMetaManager()
		userMng     = manager.NewUserManager(m)
		groupMng    = manager.NewGroupManager(m)
		relationMng = manager.NewRelationManager(m)
		msgMng      = manager.NewMsgManager(m)
		fileMng     = manager.NewFileManager(m)
	)
	if err = m.WithAdminToken(); err != nil {
		return err
	}
	ctx = m.BuildCtx(ctx)

	pro.SetContext(ctx)

	checkTasks := []*process.Task{
		process.NewTask(vars.ShouldCheckGroupNum, checker.CheckGroupNum),
		process.NewTask(vars.ShouldCheckConversationNum, checker.CheckConvNumAfterImpFriAndCrGro),
		process.NewTask(vars.ShouldCheckMessageNum, checker.CheckMessageNum),
	}

	pro.AddTasks(
		process.NewTask(vars.ShouldRegister, userMng.RegisterAllUsers),
		process.NewTask(true, userMng.InitAllSDK),
		process.NewTask(vars.IsLogin, userMng.LoginByRate),
		process.NewTask(vars.IsLogin, checker.CheckLoginByRateNum),

		process.NewTask(vars.ShouldImportFriends, relationMng.ImportFriends),
		process.NewTask(vars.ShouldImportFriends, checker.CheckLoginUsersFriends),

		process.NewTask(vars.ShouldCreateGroup, groupMng.CreateGroups),
		process.NewTask(vars.ShouldSendMsg, msgMng.SendMessages),
		//process.NewTask(vars.ShouldSendMsg, Sleep),

		process.NewTask(vars.IsLogin, userMng.LoginLastUsers),
		process.NewTask(vars.IsLogin, checker.CheckAllLoginNum),
	)

	pro.AddTasks(checkTasks...)
	pro.AddTasks(process.NewTask(vars.ShouldCheckMessageNum, statistics.MsgConsuming))

	// Uninstall and reinstall
	pro.AddTasks(
		process.NewTask(true, process.AddConditions, pro, vars.ShouldCheckUninsAndReins),
		process.NewTask(true, fileMng.DeleteLocalDB),
		process.NewTask(true, userMng.ForceLogoutAllUsers),
		process.NewTask(true, userMng.InitAllSDK),
		process.NewTask(true, userMng.LoginAllUsers),
	)
	pro.AddTasks(checkTasks...)

	return pro.Exec()
}

func main() {
	var err error
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic: %v\n", r)
			fmt.Println("Stack trace:")
			fmt.Printf("%s", debug.Stack())
			os.Exit(runFailed)
		}
		if err != nil {
			fmt.Println(utils.FormatErrorStack(err))
			os.Exit(runFailed)
		}
	}()
	ctx := context.Background()
	if err = Init(ctx); err != nil {
		log.ZError(ctx, "init err", err, "stack", utils.FormatErrorStack(err))
		fmt.Println("init err")
		return
	}
	if err = DoFlagFunc(ctx); err != nil {
		log.ZError(ctx, "do flag err", err, "stack", utils.FormatErrorStack(err))
		fmt.Println("do flag err")
		return
	}

	log.ZInfo(ctx, "start success!")
	fmt.Println("start success!")
}

func Sleep() {
	fmt.Printf("sleep %d s for sync data~\n", config.SleepSec)
	fmt.Print(formatutil.ProgressBar("Sleep", 0, config.SleepSec))
	for i := 0; i < config.SleepSec; i++ {
		fmt.Print(formatutil.ProgressBar("Sleep", i+1, config.SleepSec))
		time.Sleep(time.Second)
	}
	fmt.Print("\n")
}
