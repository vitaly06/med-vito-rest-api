import {
  Body,
  Controller,
  Get,
  Post,
  Query,
  Req,
  Res,
  UseGuards,
} from '@nestjs/common';
import { ApiOperation, ApiTags } from '@nestjs/swagger';

import type { Request, Response } from 'express';

import { AuthService } from './auth.service';
import { ChangePasswordDto } from './dto/change-password.dto';
import { ForgotPasswordDto } from './dto/forgot-password.dto';
import { SignInDto } from './dto/sign-in.dto';
import { SignUpDto } from './dto/sign-up.dto';
import { AdminSessionAuthGuard } from './guards/admin-session-auth.guard';
import { SessionAuthGuard } from './guards/session-auth.guard';

@Controller('auth')
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @Post('sign-up')
  async signUp(@Body() dto: SignUpDto) {
    return await this.authService.signUp(dto);
  }

  @Post('sign-in')
  async signIn(
    @Body() dto: SignInDto,
    @Res({ passthrough: true }) response: Response,
  ) {
    return await this.authService.signIn(dto, response);
  }

  @ApiOperation({
    summary: 'Данные сессии',
  })
  @Get('me')
  @UseGuards(SessionAuthGuard)
  async getCurrentUser(@Req() request: Request & { user: any }) {
    return await this.authService.getCurrentUser(request.user.id);
  }

  @Get('isAdmin')
  @UseGuards(AdminSessionAuthGuard)
  async isAdmin() {
    return { isAdmin: true };
  }

  @Post('logout')
  @UseGuards(SessionAuthGuard)
  async logout(
    @Req() request: Request,
    @Res({ passthrough: true }) response: Response,
  ) {
    const sessionId = request.cookies?.session_id;
    return await this.authService.logout(sessionId, response);
  }

  @ApiTags('Восстановление пароля')
  @ApiOperation({
    summary: 'Отправка кода подтверждения для восстановления пароля',
  })
  @Post('forgot-password')
  async forgotPassword(@Body() dto: ForgotPasswordDto) {
    return await this.authService.forgotPassword(dto.email);
  }

  @ApiTags('Восстановление пароля')
  @ApiOperation({
    summary: 'Проверка кода из письма для восстановления пароля',
  })
  @Post('verify-code')
  async verifyCode(@Query('code') code: string) {
    return await this.authService.verifyCode(code);
  }

  @ApiTags('Восстановление пароля')
  @ApiOperation({
    summary: 'Смена4 пароля после подтверждения кода из почты',
  })
  @Post('change-password')
  async changePassword(@Body() dto: ChangePasswordDto) {
    return await this.authService.changePassword(+dto.userId, dto.password);
  }
}
