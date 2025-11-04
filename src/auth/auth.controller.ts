import {
  Body,
  Controller,
  Post,
  Res,
  Req,
  UseGuards,
  Query,
} from '@nestjs/common';
import { AuthService } from './auth.service';
import { SignUpDto } from './dto/sign-up.dto';
import { SignInDto } from './dto/sign-in.dto';
import type { Response, Request } from 'express';
import { JwtAuthGuard } from './guards/jwt-auth.guard';
import { ForgotPasswordDto } from './dto/forgot-password.dto';
import { ApiOperation, ApiTags } from '@nestjs/swagger';
import { ChangePasswordDto } from './dto/change-password.dto';

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

  @Post('logout')
  @UseGuards(JwtAuthGuard)
  async logout(
    @Req() request: Request & { user: any },
    @Res({ passthrough: true }) response: Response,
  ) {
    return await this.authService.logout(request.user.id, response);
  }

  @Post('refresh')
  async refresh(
    @Req() request: Request,
    @Res({ passthrough: true }) response: Response,
  ) {
    const refreshToken = request.cookies?.refresh_token;
    if (!refreshToken) {
      throw new Error('Refresh токен не найден');
    }
    return await this.authService.refreshTokens(refreshToken, response);
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
