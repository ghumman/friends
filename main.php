<?php
  header("Access-Control-Allow-Origin: *");

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
                    $stmt = $conn->prepare("SELECT password, salt FROM users WHERE email = :email");
                    $stmt->bindParam(':email', $email);
                    $stmt->execute();

                    if ($stmt->rowCount() > 0) {
                        $response['status'] = 400;
                        $response['error'] = true;
                        $response['message'] = 'User Already Exists';
                        $response['time'] = date("Y-m-d H:i:s");
                    } else {


                        $salt = generateSalt();
                        $dbPassword = base64_encode(generateKey($password, $salt));

                        $stmt = $conn->prepare("select id from users order by id desc limit 1");
                        $stmt->execute();
                        $results = $stmt->fetchAll(PDO::FETCH_ASSOC);

                        $userID = 1;
                        if ($stmt->rowCount() > 0) {
                            $userID = $results[0]['id'] + 1;
                        }

                        $currentTime = date("Y-m-d H:i:s");

                        $stmt = $conn->prepare("insert into users (id, auth_type, created_at, email, first_name, last_name, password, salt) VALUES (:userID, 0, :currentTime, :email, :firstName, :lastName, :dbPassword, :salt)");
                        $stmt->bindParam(':userID', $userID);
                        $stmt->bindParam(':currentTime', $currentTime);
                        $stmt->bindParam(':email', $email);
                        $stmt->bindParam(':firstName', $firstName);
                        $stmt->bindParam(':lastName', $lastName);
                        $stmt->bindParam(':dbPassword', $dbPassword);
                        $stmt->bindParam(':salt', $salt);
        
                        if ($stmt->execute()) {

                            sendNewUserEmail($email, true, "");

                            $response['status'] = 200;
                            $response['error'] = false;
                            $response['message'] = 'User Created';
                            $response['time'] = date("Y-m-d H:i:s");
                        }

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
                    $stmt = $conn->prepare("SELECT password, salt FROM users WHERE email = :email");
                    $stmt->bindParam(':email', $email);
        
                    $stmt->execute();
        
                    if ($stmt->rowCount() > 0) {
                        $results = $stmt->fetchAll(PDO::FETCH_ASSOC);
                        
                        $key = base64_encode(generateKey($password, $results[0]['salt']));

                        if ($key == $results[0]['password']) {


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
                    $stmt = $conn->prepare("SELECT password, salt FROM users WHERE email = :email");
                    $stmt->bindParam(':email', $email);
        
                    $stmt->execute();
        
                    if ($stmt->rowCount() > 0) {
                        $results = $stmt->fetchAll(PDO::FETCH_ASSOC);
                        
                        $key = base64_encode(generateKey($password, $results[0]['salt']));

                        if ($key == $results[0]['password']) {

                            $salt = generateSalt();
                            $dbPassword = base64_encode(generateKey($newPassword, $salt));

                            $stmt = $conn->prepare("UPDATE users SET salt=:salt, password=:dbPassword where email=:email");
                            $stmt->bindParam(':email', $email);
                            $stmt->bindParam(':salt', $salt);
                            $stmt->bindParam(':dbPassword', $dbPassword);
                            $stmt->execute();


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

                    $stmt = $conn->prepare("SELECT password, salt FROM users WHERE email = :email");
                    $stmt->bindParam(':email', $email);
        
                    $stmt->execute();
        
                    if ($stmt->rowCount() > 0) {

                        $token = uniqid();
                        $stmt = $conn->prepare("UPDATE users SET reset_token=:token where email=:email");
                        $stmt->bindParam(':email', $email);
                        $stmt->bindParam(':token', $token);
            
                        $stmt->execute();


                        sendNewUserEmail($email, false, $token);
                        

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

                    $stmt = $conn->prepare("select id from users where reset_token=:token");
                    $stmt->bindParam(':token', $token);

                    $stmt->execute();
                    if ($stmt->rowCount() > 0) {
                        $results = $stmt->fetchAll(PDO::FETCH_ASSOC);

                        $userID = $results[0]['id'];

                        $salt = generateSalt();
                        $dbPassword = base64_encode(generateKey($password, $salt));

                        $stmt = $conn->prepare("UPDATE users SET salt=:salt, password=:dbPassword, reset_token=NULL where id=:userID");
                        $stmt->bindParam(':salt', $salt);
                        $stmt->bindParam(':dbPassword', $dbPassword);
                        $stmt->bindParam(':userID', $userID);
                        $stmt->execute();
                        

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
                    $stmt = $conn->prepare("SELECT password, salt FROM users WHERE email = :email");
                    $stmt->bindParam(':email', $email);
    
                    $stmt->execute();
    
                    if ($stmt->rowCount() > 0) {
                        $results = $stmt->fetchAll(PDO::FETCH_ASSOC);
                        
                        $key = base64_encode(generateKey($password, $results[0]['salt']));

                        if ($key == $results[0]['password']) {

                            $users = array();

                            $stmt = $conn->prepare("select first_name, last_name, email FROM users where email!=:email");
                            $stmt->bindParam(':email', $email);
            
                            $stmt->execute();
                            $results = $stmt->fetchAll(PDO::FETCH_ASSOC);

                            for ($x=0; $x<count($results); $x++) {
                                $users[$x] = [
                                    "firstName" => $results[$x]['first_name'], 
                                    "lastName" => $results[$x]['last_name'],
                                    "email" => $results[$x]['email']
                                ];
                            }


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
                    $stmt = $conn->prepare("SELECT password, salt, id FROM users WHERE email = :email");
                    $stmt->bindParam(':email', $messageFromEmail);
    
                    $stmt->execute();
    
                    if ($stmt->rowCount() > 0) {
                        $dataSender = $stmt->fetchAll(PDO::FETCH_ASSOC);
                        
                        $key = base64_encode(generateKey($password, $dataSender[0]['salt']));

                        if ($key == $dataSender[0]['password']) {
                            $stmt = $conn->prepare("SELECT password, salt, id FROM users WHERE email = :email");
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
                    $stmt = $conn->prepare("SELECT password, salt, id FROM users WHERE email = :email");
                    $stmt->bindParam(':email', $userEmail);
    
                    $stmt->execute();
    
                    if ($stmt->rowCount() > 0) {
                        $dataSender = $stmt->fetchAll(PDO::FETCH_ASSOC);
                        
                        $key = base64_encode(generateKey($password, $dataSender[0]['salt']));

                        if ($key == $dataSender[0]['password']) {
                            $stmt = $conn->prepare("SELECT password, salt, id FROM users WHERE email = :email");
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
