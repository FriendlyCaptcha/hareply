#!/bin/sh

systemctl --system daemon-reload || true
systemctl enable hareply || true
