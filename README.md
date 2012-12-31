daemonkit
=====
create a daemonizer
<pre>
  <code>
    import dk "daemonkit"
    import "os"
    ...
    func main() {
      daemon := dk.NewDaemonizer("/tmp")
      daemon.WatchCli(os.Args)
    }
  </code>
</pre>

cli usage
<pre>
  <code>
    usage: daemonizer {start|stop|restart} prog ...args
  </code>
</pre>
