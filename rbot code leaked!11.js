const among = require("mineflayer")
const vec3 = require("vec3")
const Item = require('prismarine-item')('1.19.3')
const bot = among.createBot({
	username: "rbot." + Math.floor(Math.random() * 9999),
	host: "kaboom.fusselig.xyz",
	port: 25565,
	hideErrors: true,
	version: "1.19.3"
})

corePos = {x: 0, y: 0, z: 0}
coreSize = [16, 16, 16]
coreIndex = 0
coreAvailable = false

function tellraw(whoAsked, message) {
	core(`tellraw ${whoAsked} ${JSON.stringify(message)}`)
}
function tellraw(message) {
	core(`tellraw @a ${JSON.stringify(message)}`)
}

var chatQueue = []

function chat(message) {
	message.split("\n").forEach((r) => {
	chatQueue.push(r.substring(0, 256))		
	})
}

function core(command) {
	var sus = [0, 0, 0]
	sus[0] += coreIndex
	while(sus[0] > 15) {
		sus[0] -= 16
		sus[1] += 1
	}
	while(sus[1] > 15) {
		sus[1] -= 16
		sus[2] += 1
	}
	try {
	bot.setCommandBlock(new vec3(corePos.x + sus[0], corePos.y + sus[1], corePos.z + sus[2]), command, {
		mode: 1,
		trackOutput: true,
		conditional: false,
		alwaysActive: true
	})	
	} catch (e) {
	
	}
	coreIndex++
	
}

function refillCore() {
	bot.chat(`/fill ${corePos.x} ${corePos.y} ${corePos.z} ${corePos.x + (coreSize[0] - 1)} ${corePos.y + (coreSize[1] - 1)} ${corePos.z + (coreSize[2] - 1)} repeating_command_block destroy`)
}

commands = {
	help: {
		name: "help",
		desc: "Show command list",
		execute: async (username, args) => {
			var sus = []
			Object.values(commands).forEach(
				e => {
					sus = sus.concat([
						{text: e.name, color: "red"},
						{text: ": ", color: "dark_gray"},
						{text: e.desc, color: "light_purple"},
						"\n"
					])
				}
			)
			sus = sus.concat([
				{text: "Made using ", color: "red"},
				{text: "ChromeOS™", color: "gray"}
			])
			tellraw(sus)
		}
	},
	rc: {
		name: "rc",
		desc: "Refill Core :tm:",
		execute: async (username, args) => {
			refillCore()
		}
	}
}
bot.on("chat", (username, message) => {
	if(message.startsWith("%")) // lets say the prefix is %
	{
		var args = message.substring(1, 4096).split(" ") // max command length 4096
		var command = args.shift()
		if(command in commands) {
			commands[command].execute()
		} else {
			bot.chat("Command not found!")
		}
		
	}
})
bot.on("login", () => {
	chat("/rank &c[bot]")
  chat("/c on")
  
  
	
})
setInterval(() => {
	coreIndex = 0
	if(chatQueue.length > 0) {
		bot.chat(chatQueue.shift())
	}
}, 60)
bot.on("message", (msg) => {
	console.log(msg.toAnsi())
})
bot.on("forcedMove", () => {
	corePos = {x: Math.floor(bot.entity.position.x), y: 0, z: Math.floor(bot.entity.position.z)}
	refillCore()
	if(!coreAvailable) {
		coreAvailable = true
		setTimeout(() => {
			tellraw([
				{text: "R-Bot", color: "red"},
				{text: " by mirkokral.", color: "gray"},
				{text: "\n- Initialized!", color: "green"}
			])	
		}, 500)
    
	}
})
bot.on("end", (msg) => {
	console.log(msg)
}) 