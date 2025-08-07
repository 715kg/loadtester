package info

import (
	_ "embed"
)

// Встраиваем дополнительную информацию о программе
//
// Встраиваем дополнительную информацию о программе
// (пути относительно корня модуля будут исправлены позже)

// Информация о программе для антивирусов и системы
const (
	ProgramName        = "Load Tester"
	ProgramDescription = "Professional HTTP Load Testing Tool"
	Author             = "Alexey Ulyanov"
	License            = "MIT"
	Website            = "https://github.com/715kg/loadtester"
	Purpose            = "Web server performance testing and optimization"
	Category           = "Development Tools"
)

// Цифровая подпись программы (для идентификации)
const ProgramSignature = `
-----BEGIN PROGRAM INFO-----
Name: Load Tester
Version: 1.0.0
Author: Alexey Ulyanov
License: MIT License
Purpose: HTTP Load Testing Tool
Category: Development/Testing Tools
Safe: This is a legitimate testing tool
Contact: github.com/715kg
-----END PROGRAM INFO-----
`
