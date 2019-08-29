global SIZE
SIZE = 80

def listSplit(li,char):
	ret = []
	tmp = []
	for item in li:
		if item == char:
			ret.append(tmp)
			tmp = []
		else:
			tmp.append(item)
	ret.append(tmp)
	return ret

def textPreprocess(file):
	io = open(file)
	lines = io.readlines()
	io.close()
	chars = []
	for line in lines:
		chars.append(list(line))
	ret = []
	for i in range(SIZE):
		tmp = []
		for j in range(SIZE):
			try:
				tmp.append(chars[i][j])
			except:
				tmp.append("")
		ret.append(tmp)
	return ret

def charMapFind(charMap, char):
	i = 0
	for li in charMap:
		j = 0
		for item in li:
			if item == char:
				return [i,j]
			j = j + 1
		i = i + 1
	return [0,0]

class TwoDPointer:
	def __init__(self,charMap,startPos,direction):
		self.charMap, self.currentPos, self.direction = charMap, startPos, direction

	def next(self):
		newx = (self.currentPos[0] + self.direction[0]) %len(self.charMap)
		newy = (self.currentPos[1] + self.direction[1]) %len(self.charMap[newx])
		self.currentPos = (newx, newy)

	def getToken(self):
		return self.charMap[self.currentPos[0]][self.currentPos[1]]

	def getCurrentPos(self):
		return self.currentPos

	def getDirection(self):
		return self.direction
	
	def setCurrentPos(self,newPos):
		self.currentPos = newPos

	def setDirection(self,newDirection):
		self.direction = newDirection

global running
global quoteMode
global slist
global pos
global local_list
global local_element
global jump_list

running = True
quoteMode = False

slist = [[],[],[]]
pos = 0
jump_list = []

charMap = textPreprocess("test")
#print("charMap",charMap)

startPos = charMapFind(charMap,"$")
pointer = TwoDPointer(charMap,startPos,[0,1])


while(running):

	pointer.next()
	token, currentPos, currentDirection = pointer.getToken(), pointer.getCurrentPos(), pointer.getDirection()
	#print("Debug token:",token,"currentPos:",currentPos,"slist:",slist[pos],"jl",jump_list)
	if token == "\"":
		quoteMode = not quoteMode
	elif token == "/":
		pointer.next()
		if pointer.getToken() != "":
			slist[pos].append(pointer.getToken())
	elif token == "\\":
		pointer.next()
		if pointer.getToken() != "":
			slist[pos].append(("\\"+pointer.getToken()).encode().decode('unicode_escape'))
	elif token == ">":
		pointer.setDirection([0,1])
	elif token == "^":
		pointer.setDirection([-1,0])
	elif token == "<":
		pointer.setDirection([0,-1])
	elif token == "v":
		pointer.setDirection([1,0])
	elif token == "?":
		if len(slist[pos]) == 0:
			pointer.setDirection([abs(currentDirection[1]),abs(currentDirection[0])])
		else:
			pointer.setDirection([-abs(currentDirection[1]),-abs(currentDirection[0])])
	elif token == "=":
		if slist[pos][-1] == slist[(pos+1)%3][-1]:
			pointer.setDirection([abs(currentDirection[1]),abs(currentDirection[0])])
		else:
			pointer.setDirection([-abs(currentDirection[1]),-abs(currentDirection[0])])
	elif token == "~":
		pos = (pos+1)%3
	elif token == "!":
		jump_list.append(currentPos)
	elif token == "#":
		pointer.setCurrentPos(jump_list.pop())
	elif token == "*":
		jump_list.pop()
	elif token == "&":
		jump_list.append(jump_list[-1])
	elif token == ".":
		print(slist[pos][-1], end = '')
	elif token == ",":
		slist[pos].append(input())
	elif token == "-":
		slist[pos].pop()
	elif token == "+":
		slist[(pos+1)%3].append(slist[pos][-1])
	elif token == "%":
		slist[pos] = []
	elif quoteMode:
		if token != "":
			slist[pos].append(token)
	elif token == "@":
		running = False


print("")
#print("\nCompleted")
