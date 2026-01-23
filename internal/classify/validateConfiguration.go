// Copyright (c) 2026 Paul Sade.
//
// This file is part of the FolderFlow project.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3,
// as published by the Free Software Foundation (see the LICENSE file).
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU General Public License for more details.

package classify

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/polocto/FolderFlow/internal/config"
)

func validateConfiguration(cfg config.Config) error {
	// --- Vérification des sources ---
	if len(cfg.SourceDirs) == 0 {
		return fmt.Errorf("aucun répertoire source configuré")
	}

	for _, src := range cfg.SourceDirs {
		if src == "" {
			return fmt.Errorf("un répertoire source est vide")
		}

		info, err := os.Stat(src)
		if os.IsNotExist(err) {
			return fmt.Errorf("répertoire source n'existe pas: %s", src)
		} else if err != nil {
			return fmt.Errorf("impossible d'accéder au répertoire source %s: %w", src, err)
		} else if !info.IsDir() {
			return fmt.Errorf("le chemin source n'est pas un répertoire: %s", src)
		}

		if cfg.Regroup != nil && src == cfg.Regroup.Path {
			return fmt.Errorf("le répertoire source est identique au chemin de regroupement: %s", src)
		}
	}

	// --- Vérification des destinations ---
	if len(cfg.DestDirs) == 0 {
		return fmt.Errorf("aucun répertoire de destination configuré")
	}

	for _, dest := range cfg.DestDirs {
		if dest.Path == "" {
			return fmt.Errorf("un répertoire de destination est vide")
		}

		info, err := os.Stat(dest.Path)
		if err == nil {
			// Existe → vérifier qu'il s'agit d'un répertoire et test d'écriture
			if !info.IsDir() {
				return fmt.Errorf("le chemin de destination n'est pas un répertoire: %s", dest.Path)
			}
			if err := testWrite(dest.Path); err != nil {
				return fmt.Errorf("pas de droit d'écriture sur le répertoire de destination %s: %w", dest.Path, err)
			}
		} else if os.IsNotExist(err) {
			// N'existe pas → vérifier qu'on peut écrire dans le premier parent existant
			parent, perr := firstExistingParent(dest.Path)
			if perr != nil {
				return fmt.Errorf("aucun dossier parent existant pour %s: %w", dest.Path, perr)
			}
			if err := testWrite(parent); err != nil {
				return fmt.Errorf("pas de droit d'écriture dans le parent existant %s pour créer %s: %w", parent, dest.Path, err)
			}
		} else {
			return fmt.Errorf("impossible d'accéder au répertoire de destination %s: %w", dest.Path, err)
		}
	}

	return nil
}

// --- Trouve le premier dossier parent existant ---
func firstExistingParent(path string) (string, error) {
	parent := path
	for {
		if parent == "" || parent == "/" {
			return "", fmt.Errorf("aucun dossier parent existant trouvé pour %s", path)
		}
		info, err := os.Stat(parent)
		if err == nil && info.IsDir() {
			return parent, nil
		} else if err != nil && !os.IsNotExist(err) {
			return "", fmt.Errorf("impossible d'accéder à %s: %w", parent, err)
		}
		parent = filepath.Dir(parent)
	}
}

// --- Teste qu'on peut écrire dans le répertoire ---
func testWrite(dir string) error {
	testFile := filepath.Join(dir, ".folderflow_write_test")
	f, err := os.Create(testFile)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Remove(testFile)
}
