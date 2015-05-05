from collections import defaultdict

def getMaxScore(girls, friends):
	points_by_girl = defaultdict(int)

	for girl, answers in girls.items():
		# Check the A "7 points if G likes naughty, dirty games"
		if answers[0] == 'Y':
			points_by_girl[girl] += 7

		# Check B: 3 points for every friend of G who likes super hero
		# action figures
		for friend in friends[girl]:
			if girls[friend][1] == 'Y':
				points_by_girl[girl] += 3

		# Check C: 6 points for every friend of a friend of G, not
		# including the friends of G and G herself, who likes men in
		# leisure suits
		friends_of_friends = set()
		for friend in friends[girl]:
			friends_of_friends |= friends[friend]
		friends_of_friends -= friends[girl]
		friends_of_friends.discard(girl)

		for friend_of_friends in friends_of_friends:
			if girls[friend_of_friends][2] == 'Y':
				points_by_girl[girl] += 6

		# Check D: 4 points if G has any friend H who likes cats, and
		# no friend of H (except perhaps G) likes cats (4 points at
		# most, not 4 for every friend).
		for friend in friends[girl]:
			if girls[friend][3] == 'Y':
				any_friend_likes_cats = False
				for friend_of_friend in friends[friend]:
					if friend_of_friend != girl and girls[friend_of_friend][3] == 'Y':
						any_friend_likes_cats = True

				if not any_friend_likes_cats:
					points_by_girl[girl] += 4
					break

		# Check E: 5 points for every girl H who likes to go shopping
		# and has no possible connection with G through a chain of
		# friends (friends, friends of friends, friends of friends of
		# friends, etc.)
		no_friends_of_friend = set(girls.keys()) - get_connected_girls_by_friends(friends, girl, set())
		no_friends_of_friend.discard(girl)
		for no_friend_of_friend in no_friends_of_friend:
			if girls[no_friend_of_friend][4] == 'Y':
				points_by_girl[girl] += 5

	return max(points_by_girl.values())

def get_connected_girls_by_friends(friends, girl, visited):
	for friend in friends[girl]:
		if friend not in visited:
			visited.add(friend)
			get_connected_girls_by_friends(friends, friend, visited)

	return visited

info = map(int, raw_input().split())
girls = {}
for i in xrange(info[0]):
	girl_info = raw_input().split()
	girls[girl_info[0]] = girl_info[1:]

friends = defaultdict(set)
for i in xrange(info[1]):
	friends_set = set(raw_input().split())
	for friend in friends_set:
		friends[friend] |= friends_set
		friends[friend].discard(friend)

print getMaxScore(girls, friends)
