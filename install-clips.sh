#!/bin/bash 

svn checkout https://svn.code.sf.net/p/clipsrules/code/ clipsrules-code
mv clipsrules-code/branches/64x/core .
rm -rf clipsrules-code
cd core
make
cd -

