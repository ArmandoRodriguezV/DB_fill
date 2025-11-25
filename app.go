package main

import (
	"context"
	"log"
	"math/rand"
	"mongo-filler/db"
	"mongo-filler/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	usersCol := database.Collection("users")
	orchardsCol := database.Collection("orchards")
	plantsCol := database.Collection("plants")
	timeUseCol := database.Collection("timeuse")

	rand.Seed(time.Now().UnixNano())

	const totalUsers = 100000
	log.Println("Generando usuarios...")

	now := time.Now()
	oneMonthAgo := now.AddDate(0, -1, 0)

	for i := 0; i < totalUsers; i++ {
		userID := primitive.NewObjectID()

		// Crear tiempo de creación del usuario (aleatorio en el último mes)
		userCreatedAt := randomTimeBetween(oneMonthAgo, now)

		user := models.User{
			ID:               userID,
			Name:             randomName(),
			Email:            randomEmail(),
			Password:         "hashedPassword123",
			OrchardsID:       []string{},
			CountOrchards:    0,
			ExperienceLevel:  rand.Intn(10),
			ProfilePhoto:     "",
			CreatedAt:        userCreatedAt,
			HistoryTimeUseID: []string{},
		}

		// --- Crear huertos ---
		orchardCount := rand.Intn(5) // 0 a 4 huertos
		for j := 0; j < orchardCount; j++ {
			orchID := primitive.NewObjectID()

			// El huerto se crea después del usuario
			orchardCreatedAt := randomTimeBetween(userCreatedAt, now)

			orchard := models.Orchard{
				ID:           orchID,
				Name:         randomOrchardName(),
				Description:  "Huerto generado automáticamente",
				PlantsId:     []string{},
				Width:        float64(rand.Intn(10) + 1),
				Height:       float64(rand.Intn(10) + 1),
				State:        true,
				CreatedAt:    orchardCreatedAt,
				UpdatedAt:    now,
				TimeOfLife:   int(now.Sub(orchardCreatedAt).Hours() / 24),
				StreakOfDays: rand.Intn(50),
				CountPlants:  0,
			}

			// Insertar huerto
			_, err := orchardsCol.InsertOne(context.Background(), orchard)
			if err != nil {
				log.Println("Error insertando huerto:", err)
			}

			user.OrchardsID = append(user.OrchardsID, orchID.Hex())
			user.CountOrchards++

			// --- Crear plantas ---
			plantCount := rand.Intn(20) // 0 a 19 plantas
			for k := 0; k < plantCount; k++ {
				plantID := primitive.NewObjectID()
				plant := models.Plant{
					ID:             plantID,
					Name:           randomPlantName(),
					Species:        "Specimen X",
					ScientificName: "Planta falsa",
					Type:           randomPlantType(),
					SunRequirement: randomSunRequirement(),
					WeeklyWatering: rand.Intn(5) + 1,
					HarvestDays:    rand.Intn(120),
					SoilType:       "mixto",
					WaterPerKg:     rand.Intn(20),
					Benefits:       []string{"decoración"},
					Size:           rand.Intn(60),
					Notes:          "",
					Tags:           []string{"random"},
				}

				_, err := plantsCol.InsertOne(context.Background(), plant)
				if err != nil {
					log.Println("Error insertando planta:", err)
				}

				orchard.PlantsId = append(orchard.PlantsId, plantID.Hex())
				orchard.CountPlants++
			}
		}

		// --- Crear timeUse ---
		timeUses := rand.Intn(5) // 0 a 4 time uses
		for t := 0; t < timeUses; t++ {
			// TimeUse debe ser después de la creación del usuario
			maxHoursAgo := int(now.Sub(userCreatedAt).Hours())
			if maxHoursAgo > 0 {
				start := now.Add(time.Duration(-rand.Intn(maxHoursAgo)) * time.Hour)
				end := start.Add(time.Duration(rand.Intn(5)+1) * time.Hour)

				timeUse := models.TimeUse{
					ID:         primitive.NewObjectID(),
					TimeInit:   start,
					TimeFinal:  end,
					TotalHours: end.Sub(start).Hours(),
				}

				_, err := timeUseCol.InsertOne(context.Background(), timeUse)
				if err != nil {
					log.Println("Error insertando timeUse:", err)
				}

				user.HistoryTimeUseID = append(user.HistoryTimeUseID, timeUse.ID.Hex())
			}
		}

		// Insertar usuario
		_, err = usersCol.InsertOne(context.Background(), user)
		if err != nil {
			log.Println("Error insertando usuario:", err)
		}

		if i%5000 == 0 {
			log.Printf("%d usuarios generados...\n", i)
		}

		if i > 0 && i%10 == 0 {
			log.Printf("[Seed] Usuarios generados: %d\n", i)
		}
	}

	log.Println("Seed completado!")
}

// randomTimeBetween genera un tiempo aleatorio entre start y end
func randomTimeBetween(start, end time.Time) time.Time {
	delta := end.Unix() - start.Unix()
	if delta <= 0 {
		return start
	}
	sec := rand.Int63n(delta)
	return start.Add(time.Duration(sec) * time.Second)
}

// --- Helpers simples ---

func randomName() string {
	names := []string{
		"Carlos", "Ana", "Pedro", "Lucía", "Miguel", "Armando", "Taku", "Marie", "Rosa", "Dan",
		"Eduardo", "Valeria", "Jorge", "Sofía", "Hugo", "Fernanda", "Adrián", "Elena", "César", "Nadia",
		"Tomás", "Julieta", "Iván", "Camila", "Ángel", "Marisol", "Esteban", "Daniela", "Fabián", "Paola",
		"Santiago", "Isabela", "Leonardo", "Rebeca", "Raúl", "Alicia", "Bruno", "Diana", "Samuel", "Regina",
		"Héctor", "Miranda", "Luis", "Arlette", "Mauricio", "Fiona", "Kevin", "Ariana", "Oscar", "Selena",
		"Joel", "Amelia", "Rafael", "Noemí", "Alonso", "Patricia", "Víctor", "Teresa", "Max", "Emma",
		"David", "Pilar", "Mario", "Clara", "Rubén", "Mara", "Erick", "Ingrid", "Guillermo", "Lina",
		"Matías", "Celeste", "Pablo", "Áurea", "Nicolás", "Norma", "Rodrigo", "Samara", "Gerardo", "Maite",
		"Diego", "Bárbara", "Julián", "Marina", "Emanuel", "Lourdes", "Alan", "Fátima", "Mateo", "Sara",
		"Cristian", "Giselle", "Emilio", "Adela", "Abel", "Renata", "Gael", "Miriam", "Omar", "Beatriz",
		"Hiroshi", "Aki", "Kenta", "Rin", "Sora", "Yuki", "Haruto", "Nana", "Riku", "Aya",
	}

	last := []string{
		"García", "Hernández", "Martínez", "López", "Ramírez",
		"Flores", "Vargas", "Torres", "Reyes", "Mendoza",
		"Castillo", "Silva", "Rojas", "Romero", "Benítez",
		"Navarro", "Peña", "Ibarra", "Santos", "Núñez",
		"Ortega", "Campos", "Vega", "Salazar", "Montoya",
		"Molina", "Soto", "Cruz", "Chávez", "Rivas",
		"Delgado", "Aguilar", "Paredes", "Valencia", "Suárez",
		"Arriaga", "Cortés", "Villalobos", "Mejía", "Miranda",
		"Ponce", "Calderón", "Bravo", "Serrano", "Cárdenas",
		"Solís", "Acosta", "Fuentes", "Velasco", "Bautista",
		"Domínguez", "Morales", "Carrillo", "Estrada", "Galindo",
		"Castañeda", "Zamora", "Padilla", "Quintero", "Sandoval",
		"Ávila", "Luna", "Rangel", "Rivero", "Correa",
		"León", "Arce", "Méndez", "Quintana", "Valdez",
		"Becerra", "Arellano", "Esquivel", "Mora", "Saucedo",
		"Montes", "Peralta", "Treviño", "Cervantes", "Coronado",
		"Izquierdo", "Palacios", "Madero", "Salgado", "Esparza",
		"Ishikawa", "Nakamura", "Kobayashi", "Yamamoto", "Watanabe",
		"Sato", "Tanaka", "Fujimoto", "Takeda", "Endo",
		"Mori", "Kato", "Shimizu", "Okada", "Inoue",
		"Suzuki", "Honda", "Abe", "Kishimoto", "Ueno",
		"Moretti", "Bianchi", "Rossi", "Conti", "Russo",
	}

	return names[rand.Intn(len(names))] + " " + last[rand.Intn(len(last))]
}

func randomEmail() string {
	domains := []string{"gmail.com", "outlook.com", "yahoo.com", "fake.com"}
	return randomString(8) + "@" + domains[rand.Intn(len(domains))]
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func randomOrchardName() string {
	n := []string{
		"Huerto Sol", "Green Patch", "Mi Jardín", "La Parcela", "EcoHuerto", "VerdeVida",
		"Huerto Aurora", "Campo Sereno", "Raíz Viva", "Bosque Dulce", "Granja Alba",
		"Huerto Encanto", "Finca Armonía", "La Pradera Clara", "EcoCampo", "Huerto Horizonte",
		"Tierra Noble", "Huerto Bonanza", "Jardín Lunar", "Huerto Nativo",
		"Sendero Verde", "Huerto Brisa", "Valle Vivo", "La Colina Verde", "Huerto Primavera",
		"Loma Fresca", "Huerto Renacer", "Campo Tranquilo", "Jardín Escondido", "Huerto Susurro",
		"Tierra Serena", "Huerto Amanecer", "Granja Silvestre", "Parcela Abundante",
		"Huerto Cascada", "Montaña Verde", "Jardín del Sol", "Huerto Oasis", "Pastizal Dorado",
		"Río Verde", "Huerto Raíces", "La Tierra Fértil", "Huerto Encina", "Finca Brillante",
		"Huerto de la Luna", "Huerto Paraíso", "Campo Dorado", "Huerto Esperanza",
		"Verde Arboleda", "Huerto Cedro", "Huerto Senda Clara", "Jardín Boreal", "Huerto Boscoso",
		"El Vergel Antiguo", "Huerto del Alba", "Campo Cosecha", "Huerto Sahara",
		"Huerto Silencio", "Huerto Grano de Oro", "Rancho Verde Claro", "Huerto Colibrí",
		"Huerto del Trébol", "La Hoja Viva", "Nido Verde", "Huerto Rocío",
		"Huerto Horizonte Azul", "Jardín Estelar", "Huerto Buena Tierra", "Huerto Lirio",
		"Huerto Cascabel", "Huerto Nogal", "Prado Sereno", "Huerto Amanecer Rojo",
		"Huerto Loma Verde", "Parcela Diamante", "Huerto Encanto Solar", "Huerto Místico",
		"EcoRancho", "Huerto del Manantial", "Huerto Fénix", "Huerto Senda Verde",
		"Huerto Valle Claro", "Huerto Yermo", "Huerto del Molino", "Verde Horizonte",
		"Huerto del Canto", "Huerto Flor Azul", "Huerto Peñasco", "Huerto del Cerezo",
		"Huerto Marea Verde", "Huerto Viento Norte", "Huerto Armonía Natural", "Huerto Dorado",
		"Jardín Maravilla", "Huerto Pétalo Blanco", "Huerto Umbra", "Huerto del Sauce",
		"Granja Tierra Clara", "Huerto Bruma", "Huerto de los Olivos", "Huerto Hiedra",
		"Huerto Terra Nova", "Huerto Encina Real", "Huerto Esmeralda", "Huerto Loma Alta",
		"Campo Esencial", "Huerto Cobre", "Huerto Aroma Verde", "Huerto Mirador",
		"Huerto Encanto del Sur", "Huerto Flor del Alba", "Huerto Luna Nueva",
		"Huerto Estrella del Campo", "Huerto Río Claro", "Huerto Piedra Viva",
	}

	return n[rand.Intn(len(n))]
}

func randomPlantName() string {
	p := []string{
		"Lavanda", "Romero", "Menta", "Albahaca", "Cactus", "Helecho",
		"Tomillo", "Salvia", "Orégano", "Cilantro", "Hierbabuena",
		"Manzanilla", "Caléndula", "Ruda", "Valeriana", "Melisa",
		"Eucalipto", "Aloe Vera", "Bambú", "Lirio", "Rosa",
		"Girasol", "Tulipán", "Orquídea", "Jazmín", "Gardenia",
		"Peonía", "Magnolia", "Hortensia", "Camelia", "Begonia",
		"Dalia", "Clavel", "Violeta Africana", "Petunia", "Geranio",
		"Azalea", "Buganvilla", "Crisantemo", "Diente de León",
		"Hibisco", "Margarita", "Loto", "Nenúfar", "Bromelia",
		"Suculenta Jade", "Aloe Arborescens", "Kalanchoe", "Agave",
		"Yucca", "Dracena", "Sansevieria", "Potus", "Monstera",
		"Ficus", "Filodendro", "Costilla de Adán", "Helecho Boston",
		"Palo de Brasil", "Zamioculca", "Croton", "Calatea",
		"Maranta", "Pilea", "Tradescantia", "Coleo", "Anturio",
		"Ciclamen", "Dieffenbachia", "Schefflera", "Hiedra Inglesa",
		"Lavanda Inglesa", "Lavanda Francesa", "Cedrón", "Stevia",
		"Hierba del Sapo", "Toronjil", "Epazote", "Mejorana",
		"Perejil", "Eneldo", "Mostaza", "Ajedrea", "Verbena",
		"Milenrama", "Equinácea", "Ginkgo", "Tejo", "Cedro Rojo",
		"Pino Blanco", "Abeto Azul", "Ciprés", "Arce Japonés",
		"Encino", "Olivo", "Haya", "Aliso", "Sauce Llorón",
		"Cerezo Japonés", "Mandarino", "Limón", "Naranjo", "Toronjo",
		"Guayabo", "Limonaria", "Hierba Mora", "Flor de Mayo",
		"Flor de Cempasúchil", "Flor de Jamaica", "Flor de Calabaza",
		"Flor de Azahar", "Amapola", "Malva", "Verbena Azul",
		"Campanilla", "Digitalis", "Borraja", "Lupino", "Aciano",
		"Salvia Roja", "Salvia Azul", "Lavanda Marina", "Helecho Cuerno de Alce",
	}

	return p[rand.Intn(len(p))]
}

func randomPlantType() models.PlantType {
	types := []models.PlantType{
		models.Ornamental,
		models.Medicinal,
		models.Alimenticia,
		models.Decorativa,
	}
	return types[rand.Intn(len(types))]
}

func randomSunRequirement() models.SunRequirementType {
	types := []models.SunRequirementType{
		models.Poca,
		models.Mucha,
	}
	return types[rand.Intn(len(types))]
}
