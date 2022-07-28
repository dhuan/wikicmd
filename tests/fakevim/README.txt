FAKEVIM

!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
This program is meant to be used only as part of wikicmd's test suite.
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

Running fakevim without any options set will do nothing special, other than
outputting the file path you issued fakevim with:

$ fakevim /path/to/some/file
fakevim:/path/to/some/file

fakevim can (programatically) edit files for you, that's where the FAKEVIM_MODE
Environment Variable comes in.

FAKEVIM_MODE=append
===================

Appends FAKEVIM_CONTENT to the edited file.

FAKEVIM_MODE=overwrite
======================

Overwrites the edited file with FAKEVIM_CONTENT.

FAKEVIM_MODE=none
=================

This is the default mode. Nothing really happens - the file fakevim is opened
with is not modified at all.
