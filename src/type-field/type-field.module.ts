import { Module } from '@nestjs/common';
import { TypeFieldService } from './type-field.service';
import { TypeFieldController } from './type-field.controller';
import { AuthModule } from 'src/auth/auth.module';

@Module({
  imports: [AuthModule],
  controllers: [TypeFieldController],
  providers: [TypeFieldService],
})
export class TypeFieldModule {}
