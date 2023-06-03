-- MySQL Script generated by MySQL Workbench
-- Sat 03 Jun 2023 11:18:21 CEST
-- Model: New Model    Version: 1.0
-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema todo
-- -----------------------------------------------------
DROP SCHEMA IF EXISTS `todo` ;

-- -----------------------------------------------------
-- Schema todo
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `todo` ;
USE `todo` ;

-- -----------------------------------------------------
-- Table `todo`.`role`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `todo`.`role` ;

CREATE TABLE IF NOT EXISTS `todo`.`role` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `description` VARCHAR(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC) VISIBLE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `todo`.`user`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `todo`.`user` ;

CREATE TABLE IF NOT EXISTS `todo`.`user` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(45) NOT NULL,
  `fk_role` INT NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `username_UNIQUE` (`username` ASC) VISIBLE,
  CONSTRAINT `fk_role_1`
    FOREIGN KEY (`id`)
    REFERENCES `todo`.`role` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `todo`.`permission`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `todo`.`permission` ;

CREATE TABLE IF NOT EXISTS `todo`.`permission` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NULL,
  `description` VARCHAR(45) NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `todo`.`role_has_permission`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `todo`.`role_has_permission` ;

CREATE TABLE IF NOT EXISTS `todo`.`role_has_permission` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `fk_role` INT NOT NULL,
  `fk_permission` INT NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_role`
    FOREIGN KEY (`id`)
    REFERENCES `todo`.`role` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_permission`
    FOREIGN KEY (`id`)
    REFERENCES `todo`.`permission` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
