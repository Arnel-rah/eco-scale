package main

import (
	"fmt"
	"log"
	"github.com/Arnel-rah/eco-scale/config"
)

func main() {
	fmt.Println("=== Eco-Scale Daemon ===")

	cfg, err := config.LoadConfig("config/policy.yaml")
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la config: %v", err)
	}

	fmt.Printf("Configuration version %s chargée avec succès !\n", cfg.Version)
	fmt.Printf("Mode alerte actuel : %v\n", cfg.AlertMode)

	fmt.Println("\nPolitiques définies :")
	for _, p := range cfg.Policies {
		fmt.Printf("- %s (Priorité: %s, CPU Max: %d%%)\n", p.ContainerName, p.Priority, p.CPULimit)
	}
}
