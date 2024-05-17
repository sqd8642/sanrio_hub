package data

import (
	"time"
)

func PopulateDatabase(models Models) error {
	for _, character := range characters {
		models.Characters.Insert(&character)
	}
	// TODO: Implement restaurants pupulation
	// TODO: Implement the relationship between restaurants and menus
	return nil
}

var characters = []Character{
	{
		Name: "Cinnamoroll",
		Debut: parseTime("2001-03-06T00:00:00Z"),
		Description: "Cinnamoroll is a white puppy with long ears, blue eyes, pink cheeks, and a plump curled tail. He is known for his fluffy appearance and sweet nature.",
		Personality: "Friendly, gentle, and always full of energy.",
		Hobbies: "Flying through the sky using his ears, napping in the sun, and making new friends.",
		Affiliations: []string{"Sanrio", "Cinnamoroll Cafe", "Little Twin Stars", "Sanrio Puroland"},
	  },
	  {
		Name: "Pompompurin",
		Debut: parseTime("1996-04-16T00:00:00Z"),
		Description: "Pompompurin is a golden retriever dog character from Sanrio, known for his trademark brown beret. He has a laid-back and friendly personality.",
		Personality: "Relaxed, easygoing, and always eager to make new friends.",
		Hobbies: "Eating pudding, napping in cozy corners, and going on adventures.",
		Affiliations: []string{"Sanrio","Pompompurin Cafe", "Hello Kitty", "Sanrio Puroland"},
	  },
	  {
		Name: "Keroppi",
		Debut: parseTime("1987-07-10T00:00:00Z"),
		Description: "Keroppi is a green frog character known for his cheerful and energetic personality. He wears a red and white striped shirt.",
		Personality: "Optimistic, adventurous, and always ready for fun.",
		Hobbies: "Swimming, playing baseball, and exploring the outdoors.",
		Affiliations: []string{"Sanrio", "Keroppi Cafe", "Hello Kitty", "Sanrio Puroland"},
	  },
	  {
		Name: "Gudetama",
		Debut: parseTime("2013-03-08T00:00:00Z"),
		Description: "Gudetama is a lazy egg character known for its unmotivated and apathetic demeanor. It often lounges around in various poses.",
		Personality: "Lazy, indifferent, and always looking for a place to nap.",
		Hobbies: "Sleeping, lounging, and contemplating the meaning of existence.",
		Affiliations: []string {"Sanrio", "Gudetama Cafe", "Hello Kitty", "Sanrio Puroland"},
	  },
	  {
		Name: "Kuromi",
		Debut: parseTime("2005-06-01T00:00:00Z"),
		Description: "Kuromi is a mischievous yet cute girl character with black fur, white hair, and a pink skull bow. She often plots pranks but has a soft spot for her friends.",
		Personality: "Sassy, adventurous, and fiercely loyal.",
		Hobbies: "Designing fashion, causing mischief, and hanging out with friends.",
		Affiliations: []string{"Sanrio","Kuromi Cafe","My Melody","Sanrio Puroland"},
	  },
	  {
		Name: "Chococat",
		Debut: parseTime("1996-07-02T00:00:00Z"),
		Description: "Chococat is a black cat character with floppy ears and big, round eyes. He has a chocolate-colored nose and is often seen with his friends.",
		Personality: "Friendly, curious, and always eager to explore.",
		Hobbies: "Reading, playing with friends, and discovering new things.",
		Affiliations: []string{"Sanrio", "Chococat Cafe", "Hello Kitty", "Sanrio Puroland"},
	},
	{
		Name: "Aggretsuko",
		Debut: parseTime("2015-04-01T00:00:00Z"),
		Description: "Aggretsuko is a red panda who works in an office job. She copes with her daily frustrations through singing death metal karaoke.",
		Personality: "Frustrated, sarcastic, and secretly passionate.",
		Hobbies: "Singing karaoke, venting through music, and spending time with friends.",
		Affiliations: []string{"Sanrio", "Aggretsuko Cafe", "Hello Kitty", "Sanrio Puroland"},
	},
	{
		Name: "Patty & Jimmy",
		Debut: parseTime("2006-03-15T00:00:00Z"),
		Description: "Patty and Jimmy are twin sheep siblings known for their matching bowties. Patty is the older sister and Jimmy is the younger brother.",
		Personality: "Spirited, playful, and always looking out for each other.",
		Hobbies: "Exploring, playing pranks, and enjoying life's adventures together.",
		Affiliations: []string{"Sanrio", "Patty & Jimmy Cafe", "Hello Kitty", "Sanrio Puroland"},
	},
	{
		Name: "Hangyodon",
		Debut: parseTime("1984-01-01T00:00:00Z"),
		Description: "Hangyodon is a blue, grumpy-looking fish character with large lips and a single tooth. Despite his appearance, he has a soft heart.",
		Personality: "Grumpy, but caring, and surprisingly sensitive.",
		Hobbies: "Fishing, swimming, and enjoying peaceful moments.",
		Affiliations: []string{"Sanrio", "Hangyodon Cafe", "Hello Kitty", "Sanrio Puroland"},
	},
}

func parseTime(timeStr string) time.Time {
    t, err := time.Parse(time.RFC3339, timeStr)
    if err != nil {
        // Handle error
    }
    return t
}