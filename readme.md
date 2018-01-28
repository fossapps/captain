# go capn'

### What's Go Captain?
The question you should be asking is Who is Go Cap'n? He's captain of your ship, he controls all the jobs you need done. From cleaning up to polishing.
Go Captain takes your command, and runs it.

### Why can't I do this myself?
Well, technically you can. But if you use Capn', he'll tell you if the job failed, he'll report to you how long is it taking, etc.
You can setup a adapter which will publish the information to all your crew members on slack.

### What else can this do?
As of now, it takes your job and your run time processor, and runs them in goroutines, for every tick (configurable), your runtime processor is called
where you can note things like how long has it been running, did it send something in channel, do we need to tell our crew that this job is taking too long?

It also supports lock provider. Imagine a situation where you want a cron to run every single minute, and if one instance of this job is already running, lock provider will make sure we don't run duplicate ones. You can remove this restriction by not implementing a lock provider for a job.

I'm working on ResultProcessor right now, which is called after the job is finished, which will take in the summary of what the job did and can do stuff with it like: reporting to your crew on slack/irc/telegram/whatever

### Why would I use this?
Cron jobs comes to mind. So many times, cron jobs ends up not running/misbehaving/failing, etc. And we pretty much don't know what happened.
With this, I can setup in a way that if a particular job takes more than 1 minute, it reports to slack, it reports the summary of job to slack, it reports if something goes wrong.
You are always in the `for {}` (See what I did there? you're always in the loop)

### How good is it?
Honestly, it's not. As of now I don't know if I can use this on production, but I'm working hard on improving this. Since I don't get to go that much (see again?), I'll try to give my free time for this. And slowly start to move forward so that this can be stable in future.
I believe together we can get this capn' sail the sea which he dreams of.

### Can I help?
Of course you can. There's can always be something that I didn't think about, some edge case, some missing tests, some feature which might be handy, etc.
I'd love to get any help I can. Be it a new bug report, a new feature request, a pull request to help me fix something. I'll use labels like "help-wanted" or "good-first-pr" for the things I'd want help with. You can grab anything you like too.
