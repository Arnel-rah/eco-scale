package main

import (
	"fmt"
	"log"

	"github.com/Arnel-rah/eco-scale/config"
	"github.com/Arnel-rah/eco-scale/docker"
)

func main() {
	fmt.Println("=== Eco-Scale Daemon ===")

	cfg, err := config.LoadConfig("config/policy.yaml")
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la config: %v", err)
	}

	fmt.Printf("Configuration version %s chargée avec succès !\n", cfg.Version)
	fmt.Printf("Mode alerte actuel : %v\n\n", cfg.AlertMode)

	fmt.Println("Politiques définies :")
	for _, p := range cfg.Policies {
		fmt.Printf("- %s (Priorité: %s, CPU Max: %d%%)\n", p.ContainerName, p.Priority, p.CPULimit)
	}
	fmt.Println()

	scanner, err := docker.NewDockerScanner()
	if err != nil {
		log.Fatalf("Erreur Docker: %v", err)
	}
	defer scanner.Close()

	activeContainers, err := scanner.ListActiveContainers()
	if err != nil {
		log.Fatalf("Erreur lors de la récupération des conteneurs: %v", err)
	}

	fmt.Println("=== Conteneurs en cours d'exécution ===")
	if len(activeContainers) == 0 {
		fmt.Println("Aucun conteneur actif trouvé.")
		return
	}

	for _, c := range activeContainers {
		fmt.Printf("ID: %s | Nom: %s | Image: %s | Statut: %s\n",
			c.ID[:10],
			c.Names[0],
			c.Image,
			c.Status,
		)
	}
}
