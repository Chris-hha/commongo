package main
/*
  多函数并发串行执行控制demo
*/
func main{
  ch1 := make(chan bool)
    ch2 := make(chan bool)
    ch3 := make(chan bool)
    ch4 := make(chan bool)
    ch5 := make(chan bool)

    wg.Add(5)
    go func(ch1 chan bool) {
      defer wg.Done()
      render.ExecuteAllComponentsEnvFunc(bhc, hosts, envParamValMap, allParamValMap)
      render.GetAllComponentLocalCfgs(bhc, hostLocalCfgCh, hosts)
      ch1 <- true
    }(ch1)

    go func(ch1, ch2 chan bool) {
      defer wg.Done()
      <-ch1
      render.ExecuteAllComponentsAsyncFunc(bhc, hosts, envParamValMap, asynParamValMap, allParamValMap)
      ch2 <- true
    }(ch1, ch2)

    go func(ch2, ch3 chan bool) {
      defer wg.Done()
      <-ch2
      render.ExecuteAllComponentsDependFunc(bhc, asynParamValMap, hosts, allParamValMap)
      render.ExecuteAllComponentsAddrFunc(bhc, hosts, allParamValMap)
      ch3 <- true
    }(ch2, ch3)

    go func(ch3, ch4 chan bool) {
      defer wg.Done()
      <-ch3
      hpm = render.MergeParams(bhc, allParamValMap, hosts)
      ch4 <- true
    }(ch3, ch4)

    go func(ch4, ch5 chan bool) {
      defer wg.Done()
      <-ch4
      renderRes := render.RenderHostConfig(bhc, hostLocalCfgCh, hpm)
      hostCompInfoCh <- &renderRes
      ch5 <- true
    }(ch4, ch5)
    <-ch5
  }
