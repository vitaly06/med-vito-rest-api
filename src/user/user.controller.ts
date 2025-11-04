import { Controller, Get, Param, Req, UseGuards, Query } from '@nestjs/common';
import { UserService } from './user.service';
import { JwtAuthGuard } from 'src/auth/guards/jwt-auth.guard';
import { Request } from 'express';
import { ApiOperation } from '@nestjs/swagger';

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
}
