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
} from '@nestjs/common';
import { UserService } from './user.service';
import { JwtAuthGuard } from 'src/auth/guards/jwt-auth.guard';
import { Request } from 'express';
import { ApiBody, ApiConsumes, ApiOperation } from '@nestjs/swagger';
import { FileInterceptor } from '@nestjs/platform-express';
import { UpdateSettingsDto } from './dto/update-settings.dto';

@Controller('user')
export class UserController {
  constructor(private readonly userService: UserService) {}

  @Get('/info')
  @UseGuards(JwtAuthGuard)
  async userInfo(@Req() req: Request & { user: any }) {
    return await this.userService.getUserInfo(req.user.id);
  }

  @ApiOperation({
    summary: 'Получение номера телефона продавца',
  })
  @UseGuards(JwtAuthGuard)
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
          type: 'boolean',
          description: 'Отвечаете ли вы на звонки?',
          example: true,
        },
        photo: {
          type: 'string',
          format: 'binary',
          description: 'Аватар пользователя',
        },
      },
    },
  })
  @UseGuards(JwtAuthGuard)
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
      photo?.filename || null,
    );
  }
}
