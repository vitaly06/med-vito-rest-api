import {
  Controller,
  Get,
  Param,
  Req,
  UseGuards,
  Query,
  Post,
  UseInterceptors,
  Body,
  UploadedFile,
  Patch,
  Put,
  Delete,
} from '@nestjs/common';
import { UserService } from './user.service';
import { SessionAuthGuard } from 'src/auth/guards/session-auth.guard';
import { Request } from 'express';
import { ApiBody, ApiConsumes, ApiOperation, ApiTags } from '@nestjs/swagger';
import { FileInterceptor } from '@nestjs/platform-express';
import { UpdateSettingsDto } from './dto/update-settings.dto';
import { UpdateUserDto } from './dto/update-user.dto';
import { ProfileType } from '@prisma/client';
import { AdminSessionAuthGuard } from 'src/auth/guards/admin-session-auth.guard';

@Controller('user')
export class UserController {
  constructor(private readonly userService: UserService) {}

  @ApiOperation({
    summary: 'Возвращает всех пользователей',
  })
  @ApiTags('Управление пользователями (админка)')
  @Get('find-all')
  async findAll() {
    return await this.userService.findAll();
  }

  @ApiOperation({
    summary: 'Обновить данные пользователя',
  })
  @ApiTags('Управление пользователями (админка)')
  @UseGuards(AdminSessionAuthGuard)
  @Patch('/:id')
  async updateUser(@Param('id') id: string, @Body() dto: UpdateUserDto) {
    return await this.userService.updateUser(+id, dto);
  }

  @ApiOperation({
    summary: 'Забанить/Разбанить пользователя пользователя',
  })
  @ApiTags('Управление пользователями (админка)')
  @UseGuards(AdminSessionAuthGuard)
  @Put('/toggle-banned/:id')
  async toggleBanned(@Param('id') id: string) {
    return await this.userService.toggleBanned(+id);
  }

  @ApiOperation({
    summary: 'Удалить пользователя',
  })
  @ApiTags('Управление пользователями (админка)')
  @UseGuards(AdminSessionAuthGuard)
  @Delete('/:id')
  async deleteUser(@Param('id') id: string) {
    return await this.userService.deleteUser(+id);
  }

  @Get('/info')
  @UseGuards(SessionAuthGuard)
  async userInfo(@Req() req: Request & { user: any }) {
    return await this.userService.getUserInfo(req.user.id);
  }

  @ApiOperation({
    summary: 'Получение номера телефона продавца',
  })
  @UseGuards(SessionAuthGuard)
  @Get('show-number/:userId')
  async showNumber(
    @Param('userId') userId: string,
    @Req() req: Request & { user: any },
  ) {
    return await this.userService.getPhoneNumber(+userId, req.user.id);
  }

  @ApiOperation({
    summary: 'Обновление настроек',
  })
  @ApiConsumes('multipart/form-data')
  @ApiBody({
    description: 'Данные для обновления настроек',
    schema: {
      type: 'object',
      properties: {
        fullName: {
          type: 'string',
          description: 'ФИО',
          example: 'Садиков Виталий Дмитриевич',
        },
        phoneNumber: {
          type: 'string',
          description: 'Номер телефона',
          example: '+79510341677',
        },
        isAnswersCall: {
          type: 'string',
          description: 'Отвечаете ли вы на звонки?',
          example: 'true',
        },
        profileType: {
          type: 'string',
          enum: ['INDIVIDUAL', 'OOO', 'IP'],
          description: 'Тип профиля',
          example: 'INDIVIDUAL',
        },
        photo: {
          type: 'string',
          format: 'binary',
          description: 'Аватар пользователя',
        },
      },
    },
  })
  @UseGuards(SessionAuthGuard)
  @UseInterceptors(FileInterceptor('photo'))
  @Patch('update-settings')
  async updateSettings(
    @Req() req: Request & { user: any },
    @Body() dto: UpdateSettingsDto,
    @UploadedFile() photo?: Express.Multer.File,
  ) {
    return await this.userService.updateSettings(
      dto,
      req.user.id,
      photo || null,
    );
  }

  @ApiOperation({
    summary: 'Получение настроек пользователя',
  })
  @UseGuards(SessionAuthGuard)
  @Get('profile-settings')
  async getProfileSettings(@Req() req: Request & { user: any }) {
    return await this.userService.getProfileSettings(req.user.id);
  }

  @ApiOperation({
    summary: 'Отправка пиьсма на почту для подтверждения почты',
  })
  @UseGuards(SessionAuthGuard)
  @Post('verify-email')
  async verifyEmail(@Req() req: Request & { user: any }) {
    return await this.userService.verifyEmail(req.user.id);
  }

  @ApiOperation({
    summary: 'Проверка кода из письма для верификации почты',
  })
  @Post('verify-code')
  async verifyCode(@Query('code') code: string) {
    return await this.userService.verifyCode(code);
  }

  @ApiOperation({
    summary: 'Установить баланс пользователю',
  })
  @UseGuards(AdminSessionAuthGuard)
  @Put('set-balance/:userId')
  async setBalance(
    @Query('balance') balance: string,
    @Param('userId') userId: string,
  ) {
    return await this.userService.setBalance(+userId, balance);
  }
}
