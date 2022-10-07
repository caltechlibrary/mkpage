---
title: Snapcraft integration
---

Snapcraft integration
=====================

After the v1.0 release of MkPage Project I've started experimenting with creating snaps. The snap model of a secure container for execution is attractive though how this limits the applications behaviors is very constraining. A minor hurdle was overcome with v1.0.3 release where I think I have a way to install a Pandoc along side MkPage so that my `mkpage` command can render content via Pandoc and it's template engine. This comes at a cost. The snap installation is larger (and suitable Pandoc may already be on the system). It also means the snap does not support all the platforms of that are build via snapcraft because the Conda/Pandoc combo isn't available (e.g. s390x). As a result I think the snap has to remain "experimental".

Another challenge of the snap is MkPage as a project include many command line tools. Only the `mkpage` command is automatically aliased. In the snap model to avoid name collisions the `mkpage` command is action `/snap/bin/mkpage.mkpage` alaised as `mkpage`. But the other commands don't get automagically aliased. That means if you want to use `frontmatter` or `mkrss` it they need to have additional `snap alias` commands run to work.  Since the appeal of snap is easy install across Linux distributions this is a problem. The snap simply is not "easy" in my opinion.

One of the trends I've seen in Unix/POSIX command line tools is a shift away from "flags" to bare words. A particularly common model seen in tools like `git` or `go` is a `APP_NAME VERB OBJECT_PARAMETERS` syntax model. If you think of the `APP_NAME VERB` as a dot pair this is an echo of Oberon Operating System's dot pairs for procedure commands. In the case of snaps this is evern more so because of the approach taken to namespace a snap. Since Go binary tend to be large adopting this model would improve the disk storage foot print  of the MkPage project as well as remove one of the current problems in the snap package process.  This would require changes to how dependencies of mkpage and would encourage some code modernization.

While I suspect the snapcraft.yaml's layout section may solve distrution of manpages and default templates I haven't had time to experiment with it. The current set of pages in the documentation directory should be reformat to allow production of man pages as well as general MkPage documentation, but doing so would need to be integrated easily into the snap build and deployment process.

