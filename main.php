<?php

    require 'vendor/autoload.php';
    header("Access-Control-Allow-Origin: *");
    $client = new MongoDB\Client("mongodb://localhost:27017");
    $userCollection = $client->friends_mongo->user;
    $messageCollection = $client->friends_mongo->message;

    require_once('db_connect.php');
    require_once "Mail.php";

    $response = array();


    switch ($_SERVER['REQUEST_URI']) {
        case '/add-user':

            if (isTheseParametersAvailable(array('email', 'password', 'firstName','lastName'))) {
                $email = $_POST['email'];
                $password = $_POST['password'];
                $firstName = $_POST['firstName'];
                $lastName = $_POST['lastName'];
        
                try {

                    $user = $userCollection->findOne([
                        'email' => $email,
                    ]);

                    if (!is_null($user)) {
                        $response['status'] = 400;
                        $response['error'] = true;
                        $response['message'] = 'User Already Exists';
                        $response['time'] = date("Y-m-d H:i:s");
                    } else {

                        $salt = generateSalt();
                        $dbPassword = base64_encode(generateKey($password, $salt));

                        $currentTime = date("Y-m-d H:i:s");

                        $insertOneResult = $userCollection->insertOne([
                            'createdAt' => $currentTime,
                            'email' => $email,
                            'firstName' => $firstName,
                            'lastName' => $lastName,
                            'dbPassword' => $dbPassword,
                            'salt' => $salt,
                        ]);

                        $response['status'] = 200;
                        $response['error'] = false;
                        $response['message'] = 'User Created';
                        $response['time'] = date("Y-m-d H:i:s");

                    }
                    
                } // try ends
        
                catch (PDOException $e) {
                    $response['status'] = 400;
                    $response['error'] = true;
                    $response['message'] = 'Error Occured, Please contact Admin, Error:'.$e;
                    $response['time'] = date("Y-m-d H:i:s");
                }	// catch ends
            } else {
                $response['status'] = 400;
                $response['error'] = true;
                $response['message'] = 'Required data missing';
                $response['time'] = date("Y-m-d H:i:s");
            }
        
        break;

        case '/login':

            if (isTheseParametersAvailable(array('email', 'password'))) {
                $email = $_POST['email'];
                $password = $_POST['password'];

                try {

                    $user = $userCollection->findOne([
                        'email' => $email,
                    ]);
        
                    if (!is_null($user)) {
                        
                        $key = base64_encode(generateKey($password, $user['salt']));

                        if ($key == $user['password']) {

                            $response['status'] = 200;
                            $response['error'] = false;
                            $response['message'] = 'Logged In';
                            $response['time'] = date("Y-m-d H:i:s");
                        } else {
                            $response['status'] = 400;
                            $response['error'] = true;
                            $response['message'] = 'Login Failed';
                            $response['time'] = date("Y-m-d H:i:s");
                        }
                    } else {
                        $response['status'] = 400;
                        $response['error'] = true;
                        $response['message'] = 'User Does Not Exist';
                        $response['time'] = date("Y-m-d H:i:s");
                    }
                }	// try ends
                catch (PDOException $e) {
                    $response['status'] = 400;
                    $response['error'] = true;
                    $response['message'] = 'Error Occured, Please contact Admin, Error:'.$e;
                    $response['time'] = date("Y-m-d H:i:s");
                }	// catch ends
            } 	// if all arguments are not present ends
            else {
                $response['status'] = 400;
                $response['error'] = true;
                $response['message'] = 'Email or Password missing';
                $response['time'] = date("Y-m-d H:i:s");
            }
        break;


        case '/change-password':

            if (isTheseParametersAvailable(array('email', 'password', 'newPassword'))) {
                $password = $_POST['password'];
                $newPassword = $_POST['newPassword'];
                $email = $_POST['email'];

                try {
                    $user = $userCollection->findOne([
                        'email' => $email,
                    ]);
        
                    if (!is_null($user)) {
                        
                        $key = base64_encode(generateKey($password, $user['salt']));

                        if ($key == $user['password']) {

                            $salt = generateSalt();
                            $dbPassword = base64_encode(generateKey($newPassword, $salt));

                            $updateResult = $userCollection->updateOne(
                                [ 'email' => $email ],
                                [ '$set' => [ 
                                    'salt' => $salt,
                                    'password' => $dbPassword
                                ]]
                            );

                            $response['status'] = 200;
                            $response['error'] = false;
                            $response['message'] = 'Password changed';
                            $response['time'] = date("Y-m-d H:i:s");
                        } else {
                            $response['status'] = 400;
                            $response['error'] = true;
                            $response['message'] = 'Original password not right';
                            $response['time'] = date("Y-m-d H:i:s");
                        }
                    } else {
                        $response['status'] = 400;
                        $response['error'] = true;
                        $response['message'] = 'User Does Not Exist';
                        $response['time'] = date("Y-m-d H:i:s");
                    }
                } catch (PDOException $e) {
                    $response['status'] = 400;
                    $response['error'] = true;
                    $response['message'] = 'Error Occured, Please contact Admin, Error:'.$e;
                    $response['time'] = date("Y-m-d H:i:s");
                } // catch ends
            } else {
                $response['status'] = 400;
                $response['error'] = true;
                $response['message'] = 'Required data missing';
                $response['time'] = date("Y-m-d H:i:s");
            }
        break;

        case '/forgot-password':

            if (isTheseParametersAvailable(array('email'))) {
                $email = $_POST['email'];
    
                try {

                    $user = $userCollection->findOne([
                        'email' => $email,
                    ]);
        
                    if (!is_null($user)) {

                        $token = uniqid();
                        $updateResult = $userCollection->updateOne(
                            [ 'email' => $email ],
                            [ '$set' => [ 
                                'resetToken' => $token
                            ]]
                        );

                        // sendNewUserEmail($email, false, $token);
                    
                        $response['status'] = 200;
                        $response['error'] = false;
                        $response['message'] = 'Reset password is sent';
                        $response['time'] = date("Y-m-d H:i:s");
                
                    } else {
                        $response['status'] = 400;
                        $response['error'] = true;
                        $response['message'] = 'User Does Not Exist';
                        $response['time'] = date("Y-m-d H:i:s");
                    }

                } catch (PDOException $e) {
                    $response['status'] = 400;
                    $response['error'] = true;
                    $response['message'] = 'Error Occured, Please contact Admin, Error:'.$e;
                    $response['time'] = date("Y-m-d H:i:s");
                } // catch ends
            } else {
                $response['status'] = 400;
                $response['error'] = true;
                $response['message'] = 'Email is required';
                $response['time'] = date("Y-m-d H:i:s");
            }
        break;

        case '/reset-password':

            if (isTheseParametersAvailable(array('token', 'password'))) {
                $token = $_POST['token'];
                $password = $_POST['password'];

                try {

                    $user = $userCollection->findOne([
                        'resetToken' => $token,
                    ]);
        
                    if (!is_null($user)) {

                        $salt = generateSalt();
                        $dbPassword = base64_encode(generateKey($password, $salt));

                        $updateResult = $userCollection->updateOne(
                            [ 'email' => $user['email'] ],
                            [ '$set' => [ 
                                'salt' => $salt,
                                'password' => $dbPassword,
                                'resetToken' => null
                            ]]
                        );
                        
                        $response['status'] = 200;
                        $response['error'] = false;
                        $response['message'] = 'Password successfully reset';
                        $response['time'] = date("Y-m-d H:i:s");

                    } else {
                        $response['status'] = 400;
                        $response['error'] = true;
                        $response['message'] = 'Token is not valid';
                        $response['time'] = date("Y-m-d H:i:s");
                    }

                } catch (PDOException $e) {
                    $response['status'] = 400;
                    $response['error'] = true;
                    $response['message'] = 'Error Occured, Please contact Admin, Error:'.$e;
                    $response['time'] = date("Y-m-d H:i:s");
                } // catch ends
            } else {
                $response['status'] = 400;
                $response['error'] = true;
                $response['message'] = 'Required data missing';
                $response['time'] = date("Y-m-d H:i:s");
            }
        break;

        case '/all-friends':

            if (isTheseParametersAvailable(array('email', 'password'))) {
                $email = $_POST['email'];
                $password = $_POST['password'];
    
                try {

                    $user = $userCollection->findOne([
                        'email' => $email,
                    ]);
        
                    if (!is_null($user)) {
                        
                        $key = base64_encode(generateKey($password, $user['salt']));

                        if ($key == $user['password']) {

                            $users = array();

                            $friends = $userCollection->find([
                                'email' => ['$ne' => $email],
                            ]);
                            
                            foreach ($friends as $friend) {
                                $info = [
                                    "firstName" => $friend['firstName'], 
                                    "lastName" => $friend['lastName'],
                                    "email" => $friend['email']
                                ];
                                array_push($users, $info);
                             };

                            $response['status'] = 200;
                            $response['error'] = false;
                            $response['message'] = 'Friends attached';
                            $response['time'] = date("Y-m-d H:i:s");
                            $response['usersAll'] = $users;
                        } else {
                            $response['status'] = 400;
                            $response['error'] = true;
                            $response['message'] = 'Login Failed';
                            $response['time'] = date("Y-m-d H:i:s");
                        }
                    } else {
                        $response['status'] = 400;
                        $response['error'] = true;
                        $response['message'] = 'User Does Not Exist';
                        $response['time'] = date("Y-m-d H:i:s");
                    }
                }	// try ends
                catch (PDOException $e) {
                    $response['status'] = 400;
                    $response['error'] = true;
                    $response['message'] = 'Error Occured, Please contact Admin, Error:'.$e;
                    $response['time'] = date("Y-m-d H:i:s");
                }	// catch ends
            } 	// if all arguments are not present ends
            else {
                $response['status'] = 400;
                $response['error'] = true;
                $response['message'] = 'Email or Password missing';
                $response['time'] = date("Y-m-d H:i:s");
            }
        break;


        case '/send-message':

            if (isTheseParametersAvailable(array('message', 'messageFromEmail', 'messageToEmail', 'password'))) {
                $message = $_POST['message'];
                $messageFromEmail = $_POST['messageFromEmail'];
                $messageToEmail = $_POST['messageToEmail'];
                $password = $_POST['password'];
    
                try {
                    $stmt = $conn->prepare("SELECT password, salt, id FROM user WHERE email = :email");
                    $stmt->bindParam(':email', $messageFromEmail);
    
                    $stmt->execute();
    
                    if ($stmt->rowCount() > 0) {
                        $dataSender = $stmt->fetchAll(PDO::FETCH_ASSOC);
                        
                        $key = base64_encode(generateKey($password, $dataSender[0]['salt']));

                        if ($key == $dataSender[0]['password']) {
                            $stmt = $conn->prepare("SELECT password, salt, id FROM user WHERE email = :email");
                            $stmt->bindParam(':email', $messageToEmail);
            
                            $stmt->execute();
            
                            if ($stmt->rowCount() > 0) {
                                $dataReceiver = $stmt->fetchAll(PDO::FETCH_ASSOC);

                                # Sender credentials are correct and both sender and receiver exists
                                $stmt = $conn->prepare("SELECT id FROM message ORDER BY id DESC LIMIT 1");
                                $stmt->execute();
                                $results = $stmt->fetchAll(PDO::FETCH_ASSOC);
                        
                                $messageID = 1;
                                if ($stmt->rowCount() > 0) {
                                    $messageID = $results[0]['id'] + 1;
                                }
                        
                                $currentTime = date("Y-m-d H:i:s");
                        
                                $stmt = $conn->prepare("insert into message (id, message, sent_at, message_from_id, message_to_id) VALUES (:messageID, :message, :currentTime, :senderID, :receiverID)");
                                $stmt->bindParam(':messageID', $messageID);
                                $stmt->bindParam(':message', $message);
                                $stmt->bindParam(':currentTime', $currentTime);
                                $stmt->bindParam(':senderID', $dataSender[0]['id']);
                                $stmt->bindParam(':receiverID', $dataReceiver[0]['id']);
                        
                                $stmt->execute();

                                $response['status'] = 200;
                                $response['error'] = false;
                                $response['message'] = 'Message sent';
                                $response['time'] = date("Y-m-d H:i:s");

                            } else {
                                $response['status'] = 400;
                                $response['error'] = true;
                                $response['message'] = 'Receiver Does Not Exist';
                                $response['time'] = date("Y-m-d H:i:s");
                            }
                        } else {
                            $response['status'] = 400;
                            $response['error'] = true;
                            $response['message'] = 'Login Failed';
                            $response['time'] = date("Y-m-d H:i:s");
                        }
                    } else {
                        $response['status'] = 400;
                        $response['error'] = true;
                        $response['message'] = 'Sender Does Not Exist';
                        $response['time'] = date("Y-m-d H:i:s");
                    }
                }	// try ends
                catch (PDOException $e) {
                    $response['status'] = 400;
                    $response['error'] = true;
                    $response['message'] = 'Error Occured, Please contact Admin, Error:'.$e;
                    $response['time'] = date("Y-m-d H:i:s");
                }	// catch ends
            } 	// if all arguments are not present ends
            else {
                $response['status'] = 400;
                $response['error'] = true;
                $response['message'] = 'Required data missing';
                $response['time'] = date("Y-m-d H:i:s");
            }
        break;

        case '/messages-user-and-friend':

            if (isTheseParametersAvailable(array('userEmail', 'friendEmail', 'password'))) {
                $userEmail = $_POST['userEmail'];
                $friendEmail = $_POST['friendEmail'];
                $password = $_POST['password'];

                
    
                try {
                    $stmt = $conn->prepare("SELECT password, salt, id FROM user WHERE email = :email");
                    $stmt->bindParam(':email', $userEmail);
    
                    $stmt->execute();
    
                    if ($stmt->rowCount() > 0) {
                        $dataSender = $stmt->fetchAll(PDO::FETCH_ASSOC);
                        
                        $key = base64_encode(generateKey($password, $dataSender[0]['salt']));

                        if ($key == $dataSender[0]['password']) {
                            $stmt = $conn->prepare("SELECT password, salt, id FROM user WHERE email = :email");
                            $stmt->bindParam(':email', $friendEmail);
            
                            $stmt->execute();
            
                            if ($stmt->rowCount() > 0) {
                                $dataReceiver = $stmt->fetchAll(PDO::FETCH_ASSOC);

                                # Sender credentials are correct and both sender and receiver exists
                                $stmt = $conn->prepare("SELECT message, message_from_id, message_to_id, sent_at  FROM message m WHERE (m.message_from_id = :senderID and m.message_to_id = :receiverID) or (m.message_from_id = :receiverID and m.message_to_id = :senderID) order by m.sent_at");
                                $stmt->bindParam(':senderID', $dataSender[0]['id']);
                                $stmt->bindParam(':receiverID', $dataReceiver[0]['id']);
                                $stmt->execute();
                                $results = $stmt->fetchAll(PDO::FETCH_ASSOC);

                                $messages = array();

                                for ($x=0; $x<count($results); $x++) {
                                    if ($results[$x]['message_from_id'] == $dataSender[0]['id'] && $results[$x]['message_to_id'] == $dataReceiver[0]['id']) {
                                        $messages[$x] = [
                                            "message" => $results[$x]['message'], 
                                            "messageFromEmail" => $userEmail,
                                            "messageToEmail" => $friendEmail,
                                            "sentAt" => $results[$x]['sent_at']
                                        ]; 
                                    } else {
                                        $messages[$x] = [
                                            "message" => $results[$x]['message'], 
                                            "messageFromEmail" => $friendEmail,
                                            "messageToEmail" => $userEmail,
                                            "sentAt" => $results[$x]['sent_at']
                                        ]; 
                                    }
                                }

                                $response['status'] = 200;
                                $response['error'] = false;
                                $response['message'] = 'Messages attached';
                                $response['time'] = date("Y-m-d H:i:s");
                                $response['usersAll'] = $messages;

                            } else {
                                $response['status'] = 400;
                                $response['error'] = true;
                                $response['message'] = 'Receiver Does Not Exist';
                                $response['time'] = date("Y-m-d H:i:s");
                            }
                        } else {
                            $response['status'] = 400;
                            $response['error'] = true;
                            $response['message'] = 'Login Failed';
                            $response['time'] = date("Y-m-d H:i:s");
                        }
                    } else {
                        $response['status'] = 400;
                        $response['error'] = true;
                        $response['message'] = 'Sender Does Not Exist';
                        $response['time'] = date("Y-m-d H:i:s");
                    }
                }	// try ends
                catch (PDOException $e) {
                    $response['status'] = 400;
                    $response['error'] = true;
                    $response['message'] = 'Error Occured, Please contact Admin, Error:'.$e;
                    $response['time'] = date("Y-m-d H:i:s");
                }	// catch ends
            } 	// if all arguments are not present ends
            else {
                $response['status'] = 400;
                $response['error'] = true;
                $response['message'] = 'Required data missing';
                $response['time'] = date("Y-m-d H:i:s");
            }
        break;

        default:
            $response['error'] = true;
            $response['message'] = 'Invalid Operation Called';
    }

    echo json_encode($response);

    function isTheseParametersAvailable($params) {
        foreach ($params as $param) {
            if (!isset($_POST[$param])) {
                return false;
            }
        }
        return true;
    }

    function generateKey($password, $salt) {
        return hash_pbkdf2("sha1", $password, $salt, 10000, 32, true);
    }

    function generateSalt() {
        $ALPHABET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
        $returnValue = "";

        for ($x = 0; $x < 30; $x++) {
            $returnValue = $returnValue.$ALPHABET[rand(0, strlen($ALPHABET) - 1)];
        }

        return $returnValue;
    }

    function saveMessage($message, $senderID, $receiverID) {
        $stmt = $conn->prepare("SELECT id FROM message ORDER BY id DESC LIMIT 1");
        $stmt->execute();
        $results = $stmt->fetchAll(PDO::FETCH_ASSOC);

        $messageID = 1;
        if ($stmt->rowCount() > 0) {
            $messageID = $results[0]['id'] + 1;
        }

        $currentTime = date("Y-m-d H:i:s");

        $stmt = $conn->prepare("insert into message (id, message, sent_at, message_from_id, message_to_id) VALUES (:messageID, :message, :currentTime, :senderID, :receiverID)");
        $stmt->bindParam(':messageID', $messageID);
        $stmt->bindParam(':message', $message);
        $stmt->bindParam(':currentTime', $currentTime);
        $stmt->bindParam(':senderID', $senderID);
        $stmt->bindParam(':receiverID', $receiverID);

        $stmt->execute();
    }

    function sendNewUserEmail($email, $accountCreated, $token) {

        $host = "ssl://smtp.gmail.com";
        $username = "username@example.com";
        $password = "password";
        $port = "465";
        $to = $email;
        $email_from = "username@example.com";
        
        

        $email_subject = $accountCreated ? "Subject: Welcome to Friends" : "Subject: Password Reset Request" ;

        $email_body = $accountCreated ? "New Account Created." :  "To reset your password, click the link below:\nhttp://localhost:8000/#/reset-password?token=".$token;

        $headers = array ('From' => $email_from, 'To' => $to, 'Subject' => $email_subject);
        $smtp = Mail::factory('smtp', array ('host' => $host, 'port' => $port, 'auth' => true, 'username' => $username, 'password' => $password));
        $mail = $smtp->send($to, $headers, $email_body);
    }
