/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`chat` /*!40100 DEFAULT CHARACTER SET latin1 */;

USE `chat`;

/*Table structure for table `messages` */

DROP TABLE IF EXISTS `messages`;

CREATE TABLE `messages` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `ownerId` char(20) NOT NULL,
  `ownerType` char(5) NOT NULL,
  `chatId` char(36) NOT NULL,
  `senderId` char(20) NOT NULL,
  `destinationId` char(20) NOT NULL,
  `destinationType` char(5) NOT NULL,
  `msg` json NOT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`ownerId`,`chatId`),
  KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=latin1;

/*Data for the table `messages` */

insert  into `messages`(`id`,`ownerId`,`ownerType`,`chatId`,`senderId`,`destinationId`,`destinationType`,`msg`,`createdAt`) values 
(13,'085246497497','user','1ec225dd-26ac-414c-8563-7af47759be02','085246497497','085246497498','user','{\"status\": \"delivered\", \"content\": \"tess\", \"client_time\": 1581256391.187, \"sender_name\": \"Taufik\", \"server_time\": 1581256391, \"content_type\": \"text\", \"delivered_time\": 1581256391}','2020-02-09 20:53:11'),
(18,'085246497497','user','7c35a01a-c524-4fd3-bf0f-16953d57248e','085246497498','085246497497','user','{\"status\": \"delivered\", \"content\": \"aaa\", \"client_time\": 1581260172.931, \"sender_name\": \"Topik\", \"server_time\": 1581260172, \"content_type\": \"text\", \"delivered_time\": 1581260172}','2020-02-09 21:56:12'),
(15,'085246497497','user','83f833c3-6b85-4861-9de6-f4fadd2bd4a2','085246497498','085246497497','user','{\"status\": \"readed\", \"content\": \"aaaa\", \"client_time\": 1581256965.848, \"readed_time\": 1581259443, \"sender_name\": \"Topik\", \"server_time\": 1581256965, \"content_type\": \"text\", \"delivered_time\": 1581256965}','2020-02-09 21:02:45'),
(14,'085246497498','user','1ec225dd-26ac-414c-8563-7af47759be02','085246497497','085246497498','user','{\"status\": \"delivered\", \"content\": \"tess\", \"client_time\": 1581256391.187, \"sender_name\": \"Taufik\", \"server_time\": 1581256391, \"content_type\": \"text\", \"delivered_time\": 1581256391}','2020-02-09 20:53:11'),
(17,'085246497498','user','7c35a01a-c524-4fd3-bf0f-16953d57248e','085246497498','085246497497','user','{\"status\": \"delivered\", \"content\": \"aaa\", \"client_time\": 1581260172.931, \"sender_name\": \"Topik\", \"server_time\": 1581260172, \"content_type\": \"text\", \"delivered_time\": 1581260172}','2020-02-09 21:56:12'),
(16,'085246497498','user','83f833c3-6b85-4861-9de6-f4fadd2bd4a2','085246497498','085246497497','user','{\"status\": \"readed\", \"content\": \"aaaa\", \"client_time\": 1581256965.848, \"readed_time\": 1581259443, \"sender_name\": \"Topik\", \"server_time\": 1581256965, \"content_type\": \"text\", \"delivered_time\": 1581256965}','2020-02-09 21:02:45');

/*Table structure for table `users` */

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` char(20) NOT NULL,
  `name` varchar(100) DEFAULT '',
  `email` varchar(50) DEFAULT '',
  `notifToken` varchar(255) DEFAULT '',
  `lastSeen` datetime DEFAULT NULL,
  `status` enum('offline','online') DEFAULT 'offline',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

/*Data for the table `users` */

insert  into `users`(`id`,`name`,`email`,`notifToken`,`lastSeen`,`status`) values 
('085246497497','Taufik','','',NULL,'offline'),
('085246497498','Topik','','',NULL,'offline'),
('085246497499','','','',NULL,'offline');

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
