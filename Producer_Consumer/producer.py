import xmlrpclib
import random
import time
proxy = xmlrpclib.ServerProxy("http://localhost:8080/")

print "The Producer connected the server"

adjectives = ["hard",
		"boundless",
		"shrill",
		"bashful",
		"opposite",
		"fluffy",
		"dear",
		"astonishing",
		"eight",
		"sick",
		"placid",
		"ad hoc",
		"ambiguous",
		"irate",
		"ordinary",
		"numerous",
		"brawny",
		"harsh",
		"calm",
		"jumbled" 
]

nouns = ["snails",
		"believe",
		"apparatus",
		"horn",
		"eggs",
		"desire",
		"snail",
		"fireman",
		"pets",
		"stocking",
		"curtain",
		"prose",
		"doctor",
		"expansion",
		"fish",
		"hammer",
		"tail",
		"profit",
		"grip",
		"regret"
];
verbs = ["lock",
		"apologise",
		"knock",
		"advise",
		"scatter",
		"nest",
		"bomb",
		"roll",
		"decorate",
		"memorise",
		"name",
		"store",
		"harass",
		"remain",
		"last",
		"hop",
		"yell",
		"mug",
		"object",
		"weigh"
];
ponctuation = [ "?", "!" ,".", "...!"];


while True:
	sentence = adjectives[random.randint(0, len(adjectives)-1)] + " "

	sentence += nouns[random.randint(0, len(nouns)-1)] + " "

	sentence += verbs[random.randint(0, len(verbs)-1)] + " "

	sentence += ponctuation[random.randint(0, len(ponctuation)-1)]
	print "Producing " + sentence
	status = proxy.produce(sentence)
	print status
	time.sleep(2)