/*
Navicat MySQL Data Transfer

Source Server         : mysql
Source Server Version : 50624
Source Host           : 127.0.0.1:3306
Source Database       : ztalk

Target Server Type    : MYSQL
Target Server Version : 50624
File Encoding         : 65001

Date: 2016-08-30 21:18:50
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for classroominfo
-- ----------------------------
DROP TABLE IF EXISTS `classroominfo`;
CREATE TABLE `classroominfo` (
  `ClassRoomId` int(36) NOT NULL AUTO_INCREMENT COMMENT '??Id',
  `ClassRoomIMId` int(36) NOT NULL COMMENT '????IM?Id',
  `CreatorUserId` int(36) NOT NULL COMMENT '?????Id',
  `ClassRoomStatus` int(11) unsigned zerofill NOT NULL COMMENT '??????0 ???, 1 ????, 2 ????, 3 ????',
  `SettingStatus` blob COMMENT '?????? 1 ????????2 ???????3 ???????4 ??????',
  `CreateTime` bigint(36) NOT NULL COMMENT '??????',
  `ClassRoomName` varchar(255) DEFAULT NULL COMMENT '????',
  `ClassRoomLogo` varchar(255) DEFAULT NULL COMMENT '????',
  `Description` varchar(255) DEFAULT NULL COMMENT '????',
  `ClassRoomCourse` varchar(255) DEFAULT NULL COMMENT '????',
  `MemberList` blob COMMENT '????Id??(Json??Blob)',
  `OnLineMemberList` blob COMMENT '????Id??(Json??Blob)',
  `HandMemberList` blob COMMENT '????Id(Json??Blob)',
  `ForbidSayMemberList` blob COMMENT '??????Id(Json??Blob)',
  `SayingMemberList` blob COMMENT '??????Id(Json??Blob)',
  `ForbidHandStatus` int(36) DEFAULT NULL COMMENT '教室的禁止举手状态 0可举手 1禁止举手',
  PRIMARY KEY (`ClassRoomId`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of classroominfo
-- ----------------------------
INSERT INTO `classroominfo` VALUES ('1', '10010', '1', '00000000000', 0x0000000100000001, '1471588949', 'ZTALK-001', 'ZTALK', 'TheStart', 'Practice', 0x00000000, 0x000000050000000700000008000000090000000A00000006, 0x00000000, 0x00000000, 0x00000000, '1');

-- ----------------------------
-- Table structure for subresource
-- ----------------------------
DROP TABLE IF EXISTS `subresource`;
CREATE TABLE `subresource` (
  `SubFileId` int(11) NOT NULL AUTO_INCREMENT COMMENT '子文件Id',
  `UserId` int(11) NOT NULL,
  `FileId` int(11) NOT NULL COMMENT '文件Id',
  `SubFileName` varchar(255) NOT NULL COMMENT '子文件名称',
  `SubFilePath` varchar(255) NOT NULL COMMENT '子文件路径',
  `SubFileTime` bigint(20) NOT NULL COMMENT '插入数据库时间',
  `SubFileType` int(11) NOT NULL COMMENT '子文件类型 1 图片',
  `IsDelete` int(11) NOT NULL COMMENT '0 未被被删除 ， 1 已被删除',
  PRIMARY KEY (`SubFileId`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of subresource
-- ----------------------------
INSERT INTO `subresource` VALUES ('1', '1', '1', 'subFileName1', '/file/subFileName1', '1', '1472524977', '1');

-- ----------------------------
-- Table structure for uploadresource
-- ----------------------------
DROP TABLE IF EXISTS `uploadresource`;
CREATE TABLE `uploadresource` (
  `FileId` int(36) NOT NULL AUTO_INCREMENT COMMENT '??Id',
  `UserId` int(36) NOT NULL COMMENT '??Id',
  `FileName` varchar(255) NOT NULL COMMENT '???????',
  `FilePath` varchar(255) NOT NULL COMMENT '???????',
  `FileThumbPath` varchar(255) DEFAULT NULL COMMENT '??????????',
  `FileTime` bigint(36) NOT NULL COMMENT '???????',
  `IsDelete` int(36) DEFAULT NULL COMMENT '????? 0 ?????(??) ?1 ??? ',
  `FileType` int(11) DEFAULT NULL,
  PRIMARY KEY (`FileId`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of uploadresource
-- ----------------------------
INSERT INTO `uploadresource` VALUES ('1', '1', 'test.txt', 'files\\8732387165bf8e39\\91fb6f39327d8ccb\\38f07f89ddbebc2a\\10\\txt\\1472524977314458900.txt', '', '1472524977', '0', null);
INSERT INTO `uploadresource` VALUES ('2', '1', 'test.txt', 'files\\8732387165bf8e39\\91fb6f39327d8ccb\\38f07f89ddbebc2a\\10\\txt\\1472525837482371200.txt', '', '1472525837', '0', null);
INSERT INTO `uploadresource` VALUES ('3', '1', 'test.txt', 'C:\\Users8732387165bf8e39\\91fb6f39327d8ccb\\38f07f89ddbebc2a\\10\\txt\\1472525987497931400.txt', '', '1472525987', '0', null);
INSERT INTO `uploadresource` VALUES ('4', '1', 'test.txt', 'files\\8732387165bf8e39\\91fb6f39327d8ccb\\38f07f89ddbebc2a\\14\\txt\\1472537575965023700.txt', '', '1472537575', '0', null);
INSERT INTO `uploadresource` VALUES ('5', '2', 'test.txt', 'files\\8732387165bf8e39\\91fb6f39327d8ccb\\38f07f89ddbebc2a\\14\\txt\\1472537605287769100.txt', '', '1472537605', '0', null);
INSERT INTO `uploadresource` VALUES ('6', '2', 'test.txt', 'files\\2\\txt\\1472538331496881700.txt', '', '1472538331', '0', null);
INSERT INTO `uploadresource` VALUES ('7', '2', 'test.txt', 'files\\2\\txt\\test.txt.txt', '', '1472539373', '0', null);
INSERT INTO `uploadresource` VALUES ('8', '2', 'test.txt', 'files\\2\\txt\\temp.txt', '', '1472544933', '0', null);

-- ----------------------------
-- Table structure for userinfo
-- ----------------------------
DROP TABLE IF EXISTS `userinfo`;
CREATE TABLE `userinfo` (
  `UserId` int(36) NOT NULL AUTO_INCREMENT COMMENT '??Id',
  `ChatId` int(36) NOT NULL COMMENT '?????IM??userId',
  `Role` int(11) DEFAULT NULL COMMENT '??: 1=VIP?? 2=???? 3=???? 4=?? 6??',
  `DeviceType` int(36) DEFAULT NULL COMMENT '??????',
  `LoginName` varchar(255) DEFAULT NULL COMMENT '????',
  `UserName` varchar(255) DEFAULT NULL COMMENT '????',
  `UserIcon` varchar(255) DEFAULT NULL COMMENT '????',
  `Password` varchar(255) NOT NULL COMMENT '????',
  `YYToken` varchar(255) DEFAULT NULL COMMENT 'YY???Token',
  `ClassRoomList` blob COMMENT '????(Json??Blob)',
  PRIMARY KEY (`UserId`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of userinfo
-- ----------------------------
INSERT INTO `userinfo` VALUES ('1', '10033', '6', '1', '10001', 'jfl', 'icon', '202cb962ac59075b964b07152d234b70', null, null);
INSERT INTO `userinfo` VALUES ('2', '10034', '6', '1', '10002', 'bkj', 'icon', '202cb962ac59075b964b07152d234b70', null, null);
INSERT INTO `userinfo` VALUES ('3', '10035', '1', '1', '10003', 'jk', 'icon', '202cb962ac59075b964b07152d234b70', null, null);
INSERT INTO `userinfo` VALUES ('4', '10036', '1', '1', '10004', 'jlk', 'icon', '202cb962ac59075b964b07152d234b70', null, null);
INSERT INTO `userinfo` VALUES ('5', '10037', '1', '1', '10005', 'cgh', 'icon', '202cb962ac59075b964b07152d234b70', null, null);
INSERT INTO `userinfo` VALUES ('6', '10038', '1', '1', '10006', 'vhj', 'icon', '202cb962ac59075b964b07152d234b70', null, null);
INSERT INTO `userinfo` VALUES ('7', '10039', '1', '1', '10007', 'nlk', 'icon', '202cb962ac59075b964b07152d234b70', null, null);
INSERT INTO `userinfo` VALUES ('8', '10040', '1', '1', '10008', 'bjk', 'icon', '202cb962ac59075b964b07152d234b70', null, null);
INSERT INTO `userinfo` VALUES ('9', '10041', '1', '1', '10009', 'ghj', 'icon', '202cb962ac59075b964b07152d234b70', null, null);
INSERT INTO `userinfo` VALUES ('10', '10042', '1', '1', '10010', 'bhl', 'icon', '202cb962ac59075b964b07152d234b70', null, null);
INSERT INTO `userinfo` VALUES ('11', '10043', '2', '1', '10011', 'lucy', 'icon', '202cb962ac59075b964b07152d234b70', null, null);
