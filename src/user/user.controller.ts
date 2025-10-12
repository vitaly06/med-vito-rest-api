import { Controller, Get, Req, UseGuards } from '@nestjs/common';
import { UserService } from './user.service';
import { JwtAuthGuard } from 'src/auth/guards/jwt-auth.guard';
import { Request } from 'express';

@Controller('user')
export class UserController {
  constructor(private readonly userService: UserService) {}

  @Get('/info')
  @UseGuards(JwtAuthGuard)
  async userInfo(@Req() req: Request & { user: any }) {
    return await this.userService.getUserInfo(req.user.id);
  }

  @Get('/my-projects')
  @UseGuards(JwtAuthGuard)
  async getMyProjects(@Req() req: Request & { user: any }) {}
}
