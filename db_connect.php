<?php

try {

    //creating a new connection object using pdo
    $conn = new PDO('pgsql:host=localhost port=5432 user=postgres dbname=friends_psql password=postgres');
    // set the PDO error mode to exception
    $conn->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
} catch (PDOException $e) {
    $conn = null;
}
