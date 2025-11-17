#!/bin/bash

eget DavidHoenisch/Oxidation9 --to $HOME/.local/bin

if test -a $HOME/.local/bin/oxidation7; then
  oxidation9 bootstrap
fi
