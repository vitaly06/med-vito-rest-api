import { Module } from '@nestjs/common';
import { HttpModule } from '@nestjs/axios';
import { AddressService } from './address.service';
import { AddressController } from './address.controller';

@Module({
  imports: [HttpModule],
  controllers: [AddressController],
  providers: [AddressService],
})
export class AddressModule {}
