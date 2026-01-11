import { Module } from '@nestjs/common';

import { AuthModule } from 'src/auth/auth.module';

import { TypeFieldController } from './type-field.controller';
import { TypeFieldService } from './type-field.service';

@Module({
  imports: [AuthModule],
  controllers: [TypeFieldController],
  providers: [TypeFieldService],
})
export class TypeFieldModule {}
